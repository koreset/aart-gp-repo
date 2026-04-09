/**
 * Composable for generating Group Risk Quotation documents in Word (.docx) format.
 *
 * Mirrors the PDF output produced by generatePDF() in QuoteOutput.vue.
 * Uses the npm `docx` package and shared helpers from quoteDataHelpers.ts.
 */

import { ref } from 'vue'
import {
  Document,
  Packer,
  Paragraph,
  TextRun,
  Table,
  TableRow,
  TableCell,
  ImageRun,
  Header,
  Footer,
  AlignmentType,
  PageOrientation,
  BorderStyle,
  WidthType,
  ShadingType,
  PageNumber,
  TabStopType,
  TabStopPosition,
  VerticalAlign
} from 'docx'
import { saveAs } from 'file-saver'
import formatDateString from '@/renderer/utils/helpers.js'
import type { BenefitTitles } from '@/renderer/types/docxQuote'
import {
  calculateQuoteTotals,
  hasAnyNonFuneralBenefits,
  categoryHasNonFuneralBenefits,
  buildInitialInfoRows,
  buildPremiumSummaryRows,
  buildGroupFuneralRows,
  buildPremiumBreakdownRows,
  buildGroupFuneralBreakdownRows,
  buildCategoryCommonBenefitRows,
  buildBenefitDefinitionRows,
  buildFuneralCoverageRows,
  buildEducatorBenefitRows
} from '@/renderer/utils/quoteDataHelpers'

// ---------------------------------------------------------------------------
// Design constants — matching the PDF colour palette
// ---------------------------------------------------------------------------

const COLORS = {
  primary: '34495E',
  secondary: '34495E',
  accent: 'E74C3C',
  lightFill: 'ECF0F1',
  dark: '2C3E50',
  altRow: 'FAFAFA',
  white: 'FFFFFF',
  navy: '1E3A5F',
  orange: 'D47600',
  lightOrange: 'FFF8F0',
  lightBlue: 'F1F8FF',
  mediumGray: '586069',
  inputLine: 'C8C8C8'
}

const FONT = 'Arial'

// A4 portrait in DXA (1440 DXA = 1 inch, 567 DXA = 1 cm)
const A4 = { width: 11906, height: 16838 }
const MARGINS = { top: 850, bottom: 850, left: 1020, right: 1020 }
const CONTENT_WIDTH = A4.width - MARGINS.left - MARGINS.right // ~9866 DXA
const LANDSCAPE_CONTENT_WIDTH = A4.height - MARGINS.left - MARGINS.right // ~14798 DXA

// Half-point sizes (multiply pt by 2)
const SIZE = {
  title: 32, // 16pt
  heading: 28, // 14pt
  subheading: 24, // 12pt
  body: 20, // 10pt
  caption: 18, // 9pt
  small: 16 // 8pt
}

// Reusable border definitions
const thinBorder = { style: BorderStyle.SINGLE, size: 4, color: 'D5D8DC' }
const thinBorders = {
  top: thinBorder,
  bottom: thinBorder,
  left: thinBorder,
  right: thinBorder
}
const noBorder = { style: BorderStyle.NONE, size: 0, color: 'FFFFFF' }
const noBorders = {
  top: noBorder,
  bottom: noBorder,
  left: noBorder,
  right: noBorder
}
const cellMargins = { top: 60, bottom: 60, left: 100, right: 100 }

// ---------------------------------------------------------------------------
// Shared header & footer builder
// ---------------------------------------------------------------------------

function makeHeader(schemeName: string) {
  return new Header({
    children: [
      new Paragraph({
        border: {
          bottom: {
            style: BorderStyle.SINGLE,
            size: 6,
            color: COLORS.lightFill,
            space: 1
          }
        },
        spacing: { after: 100 },
        children: [
          new TextRun({
            text: `${schemeName} - Quotation`,
            font: FONT,
            size: SIZE.caption,
            color: COLORS.secondary
          })
        ]
      })
    ]
  })
}

function makeFooter() {
  return new Footer({
    children: [
      new Paragraph({
        border: {
          top: {
            style: BorderStyle.SINGLE,
            size: 6,
            color: COLORS.lightFill,
            space: 1
          }
        },
        tabStops: [{ type: TabStopType.RIGHT, position: TabStopPosition.MAX }],
        children: [
          new TextRun({
            text: `Generated on ${formatDateString(new Date(), true, true, true)}`,
            font: FONT,
            size: SIZE.caption,
            color: COLORS.secondary
          }),
          new TextRun({
            text: '\tPage ',
            font: FONT,
            size: SIZE.caption,
            color: COLORS.secondary
          }),
          new TextRun({
            children: [PageNumber.CURRENT],
            font: FONT,
            size: SIZE.caption,
            color: COLORS.secondary
          })
        ]
      })
    ]
  })
}

// ---------------------------------------------------------------------------
// Reusable element builders
// ---------------------------------------------------------------------------

/** Styled section heading with underline accent. */
function sectionHeading(text: string): Paragraph {
  return new Paragraph({
    spacing: { before: 200, after: 120 },
    border: {
      bottom: {
        style: BorderStyle.SINGLE,
        size: 8,
        color: COLORS.primary,
        space: 2
      }
    },
    children: [
      new TextRun({
        text,
        font: FONT,
        size: SIZE.heading,
        bold: true,
        color: COLORS.primary
      })
    ]
  })
}

/** Sub-heading for category names. */
function categoryHeading(text: string): Paragraph {
  return new Paragraph({
    spacing: { before: 160, after: 80 },
    children: [
      new TextRun({
        text,
        font: FONT,
        size: SIZE.subheading,
        bold: true,
        color: COLORS.dark
      })
    ]
  })
}

/** Simple body paragraph. */
function bodyText(
  text: string,
  options?: {
    alignment?: (typeof AlignmentType)[keyof typeof AlignmentType]
    italic?: boolean
  }
): Paragraph {
  return new Paragraph({
    alignment: options?.alignment ?? AlignmentType.JUSTIFIED,
    spacing: { after: 100 },
    children: [
      new TextRun({
        text,
        font: FONT,
        size: SIZE.body,
        color: COLORS.dark,
        italics: options?.italic ?? false
      })
    ]
  })
}

/** Create a standard header row for a data table. */
function headerRow(
  labels: string[],
  colWidths: number[],
  contentWidth: number
): TableRow {
  return new TableRow({
    tableHeader: true,
    children: labels.map(
      (label, i) =>
        new TableCell({
          width: { size: colWidths[i], type: WidthType.DXA },
          shading: { fill: COLORS.primary, type: ShadingType.CLEAR },
          borders: thinBorders,
          margins: cellMargins,
          children: [
            new Paragraph({
              alignment: i === 0 ? AlignmentType.LEFT : AlignmentType.CENTER,
              children: [
                new TextRun({
                  text: label,
                  font: FONT,
                  size: SIZE.body,
                  bold: true,
                  color: COLORS.white
                })
              ]
            })
          ]
        })
    )
  })
}

/** Create a data row with optional bold and shading. */
function dataRow(
  values: string[],
  colWidths: number[],
  options?: {
    alignments?: (typeof AlignmentType)[keyof typeof AlignmentType][]
    bold?: boolean
    fillColor?: string
  }
): TableRow {
  return new TableRow({
    children: values.map(
      (val, i) =>
        new TableCell({
          width: { size: colWidths[i], type: WidthType.DXA },
          shading: options?.fillColor
            ? { fill: options.fillColor, type: ShadingType.CLEAR }
            : undefined,
          borders: thinBorders,
          margins: cellMargins,
          children: [
            new Paragraph({
              alignment:
                options?.alignments?.[i] ??
                (i === 0 ? AlignmentType.LEFT : AlignmentType.RIGHT),
              children: [
                new TextRun({
                  text: val,
                  font: FONT,
                  size: SIZE.caption,
                  bold: options?.bold ?? i === 0,
                  color: COLORS.dark
                })
              ]
            })
          ]
        })
    )
  })
}

/** Key-value pair table (2-column, label on grey). */
function keyValueTable(
  rows: { label: string; value: string }[],
  contentWidth: number
): Table {
  const labelWidth = Math.round(contentWidth * 0.37)
  const valueWidth = contentWidth - labelWidth

  return new Table({
    width: { size: contentWidth, type: WidthType.DXA },
    columnWidths: [labelWidth, valueWidth],
    rows: rows.map(
      (row, i) =>
        new TableRow({
          children: [
            new TableCell({
              width: { size: labelWidth, type: WidthType.DXA },
              shading: { fill: COLORS.lightFill, type: ShadingType.CLEAR },
              borders: thinBorders,
              margins: cellMargins,
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: row.label,
                      font: FONT,
                      size: SIZE.body,
                      bold: true,
                      color: COLORS.dark
                    })
                  ]
                })
              ]
            }),
            new TableCell({
              width: { size: valueWidth, type: WidthType.DXA },
              shading:
                i % 2 === 1
                  ? { fill: COLORS.altRow, type: ShadingType.CLEAR }
                  : undefined,
              borders: thinBorders,
              margins: cellMargins,
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: row.value,
                      font: FONT,
                      size: SIZE.body,
                      color: COLORS.secondary
                    })
                  ]
                })
              ]
            })
          ]
        })
    )
  })
}

// ---------------------------------------------------------------------------
// Section builders
// ---------------------------------------------------------------------------

/**
 * Sections 1–3: Cover page header (logo + insurer), title, intro text,
 * and quote summary table. Portrait orientation.
 */
function buildCoverAndSummarySection(
  quote: any,
  resultSummaries: any[],
  insurer: any
) {
  const totals = calculateQuoteTotals(resultSummaries)
  const summaryRows = buildInitialInfoRows(quote, totals)

  // Build the header layout: borderless 2-column table (text left, logo right)
  const headerLeftChildren: Paragraph[] = [
    new Paragraph({
      spacing: { before: 60, after: 20 },
      children: [
        new TextRun({
          text: insurer.name,
          font: FONT,
          size: SIZE.subheading,
          bold: true,
          color: COLORS.dark
        })
      ]
    }),
    new Paragraph({
      spacing: { after: 10 },
      children: [
        new TextRun({
          text: `${insurer.address_line_1}, ${insurer.address_line_2}`,
          font: FONT,
          size: SIZE.caption,
          color: COLORS.secondary
        })
      ]
    })
  ]

  if (insurer.address_line_3?.trim()) {
    headerLeftChildren.push(
      new Paragraph({
        spacing: { after: 10 },
        children: [
          new TextRun({
            text: insurer.address_line_3,
            font: FONT,
            size: SIZE.caption,
            color: COLORS.secondary
          })
        ]
      })
    )
  }

  headerLeftChildren.push(
    new Paragraph({
      spacing: { after: 10 },
      children: [
        new TextRun({
          text: `${insurer.city}, ${insurer.province}, ${insurer.post_code}`,
          font: FONT,
          size: SIZE.caption,
          color: COLORS.secondary
        })
      ]
    }),
    new Paragraph({
      spacing: { after: 10 },
      children: [
        new TextRun({
          text: `Tel: ${insurer.telephone}`,
          font: FONT,
          size: SIZE.caption,
          color: COLORS.secondary
        })
      ]
    }),
    new Paragraph({
      spacing: { after: 10 },
      children: [
        new TextRun({
          text: `Email: ${insurer.email}`,
          font: FONT,
          size: SIZE.caption,
          color: COLORS.secondary
        })
      ]
    })
  )

  // Logo cell content
  const logoRightChildren: Paragraph[] = []
  if (insurer.logo) {
    try {
      const logoFormat = (insurer.logo_mime_type || 'image/png').split(
        '/'
      )[1] as 'png' | 'jpg' | 'jpeg' | 'gif' | 'bmp'
      const logoBuffer = Uint8Array.from(atob(insurer.logo), (c) =>
        c.charCodeAt(0)
      )
      logoRightChildren.push(
        new Paragraph({
          alignment: AlignmentType.RIGHT,
          spacing: { before: 60 },
          children: [
            new ImageRun({
              type: logoFormat === 'jpeg' ? 'jpg' : logoFormat,
              data: logoBuffer,
              transformation: { width: 210, height: 80 },
              altText: {
                title: 'Logo',
                description: 'Insurer logo',
                name: 'insurer-logo'
              }
            })
          ]
        })
      )
    } catch (e) {
      console.warn('Failed to decode logo for DOCX:', e)
      logoRightChildren.push(
        new Paragraph({
          alignment: AlignmentType.RIGHT,
          children: [
            new TextRun({
              text: '[Company Logo]',
              font: FONT,
              size: SIZE.caption,
              italics: true,
              color: 'B4B4B4'
            })
          ]
        })
      )
    }
  }

  const leftWidth = Math.round(CONTENT_WIDTH * 0.6)
  const rightWidth = CONTENT_WIDTH - leftWidth

  const headerTable = new Table({
    width: { size: CONTENT_WIDTH, type: WidthType.DXA },
    columnWidths: [leftWidth, rightWidth],
    rows: [
      new TableRow({
        children: [
          new TableCell({
            width: { size: leftWidth, type: WidthType.DXA },
            shading: { fill: COLORS.lightFill, type: ShadingType.CLEAR },
            borders: noBorders,
            margins: { top: 80, bottom: 80, left: 120, right: 60 },
            children: headerLeftChildren
          }),
          new TableCell({
            width: { size: rightWidth, type: WidthType.DXA },
            shading: { fill: COLORS.lightFill, type: ShadingType.CLEAR },
            borders: noBorders,
            margins: { top: 80, bottom: 80, left: 60, right: 120 },
            verticalAlign: VerticalAlign.CENTER,
            children:
              logoRightChildren.length > 0
                ? logoRightChildren
                : [new Paragraph({ children: [] })]
          })
        ]
      })
    ]
  })

  // Intro text
  const introText =
    insurer.introductory_text?.trim() ||
    'We are pleased to submit for your consideration the quotation you requested for the above scheme.'

  return {
    properties: {
      page: {
        size: A4,
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name) },
    footers: { default: makeFooter() },
    children: [
      headerTable,
      // Spacer
      new Paragraph({ spacing: { before: 200, after: 100 }, children: [] }),
      // Title
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { before: 100, after: 160 },
        children: [
          new TextRun({
            text: 'Group Risk Quotation',
            font: FONT,
            size: SIZE.title,
            bold: true,
            color: COLORS.dark
          })
        ]
      }),
      // Intro
      bodyText(introText),
      // Spacer
      new Paragraph({ spacing: { before: 80 }, children: [] }),
      // Quote summary table
      keyValueTable(summaryRows, CONTENT_WIDTH)
    ]
  }
}

/**
 * Section 4: Premium Summary and Group Funeral tables.
 * Landscape orientation. Conditionally includes Premium Summary
 * only if non-funeral benefits exist.
 */
function buildPremiumSummarySection(quote: any, resultSummaries: any[]) {
  const hasNonFuneral = hasAnyNonFuneralBenefits(resultSummaries)
  const children: (Paragraph | Table)[] = []

  // Premium Summary table (conditional)
  if (hasNonFuneral) {
    children.push(
      bodyText(
        'The following table provides a summary of the benefits and premiums for the scheme:'
      ),
      sectionHeading('Premium Summary')
    )

    const premiumRows = buildPremiumSummaryRows(resultSummaries)
    const cols = [
      'Category',
      'No of Lives',
      'Total Salary',
      'Total Sum Assured',
      'Annual Premium',
      '% Salary'
    ]
    // Distribute columns across landscape width
    const cw = LANDSCAPE_CONTENT_WIDTH
    const colWidths = [
      Math.round(cw * 0.2),
      Math.round(cw * 0.12),
      Math.round(cw * 0.2),
      Math.round(cw * 0.2),
      Math.round(cw * 0.18),
      Math.round(cw * 0.1)
    ]
    // Adjust last column to absorb rounding
    colWidths[5] = cw - colWidths.slice(0, 5).reduce((a, b) => a + b, 0)

    const alignments = [
      AlignmentType.LEFT,
      AlignmentType.CENTER,
      AlignmentType.RIGHT,
      AlignmentType.RIGHT,
      AlignmentType.RIGHT,
      AlignmentType.CENTER
    ]

    const tableRows = [headerRow(cols, colWidths, cw)]
    premiumRows.forEach((row, i) => {
      const isTotal = row.category === 'Total'
      tableRows.push(
        dataRow(
          [
            row.category,
            row.memberCount,
            row.totalSalary,
            row.totalSumAssured,
            row.annualPremium,
            row.percentSalary
          ],
          colWidths,
          {
            alignments,
            bold: isTotal || undefined,
            fillColor: isTotal
              ? COLORS.lightFill
              : i % 2 === 1
                ? 'FCFDFE'
                : undefined
          }
        )
      )
    })

    children.push(
      new Table({
        width: { size: cw, type: WidthType.DXA },
        columnWidths: colWidths,
        rows: tableRows
      })
    )
  }

  // Group Funeral table (always shown)
  children.push(
    new Paragraph({ spacing: { before: 200 }, children: [] }),
    sectionHeading('Group Funeral')
  )

  const funeralRows = buildGroupFuneralRows(resultSummaries)
  const fCols = [
    'Category',
    'No of Lives',
    'Monthly Premium',
    'Annual Premium',
    'Total Annual Premium'
  ]
  const funeralCw = LANDSCAPE_CONTENT_WIDTH
  const fColWidths = [
    Math.round(funeralCw * 0.22),
    Math.round(funeralCw * 0.14),
    Math.round(funeralCw * 0.2),
    Math.round(funeralCw * 0.2),
    Math.round(funeralCw * 0.24)
  ]
  fColWidths[4] = funeralCw - fColWidths.slice(0, 4).reduce((a, b) => a + b, 0)

  const fAlignments = [
    AlignmentType.LEFT,
    AlignmentType.CENTER,
    AlignmentType.RIGHT,
    AlignmentType.RIGHT,
    AlignmentType.RIGHT
  ]

  const fTableRows = [headerRow(fCols, fColWidths, funeralCw)]
  funeralRows.forEach((row, i) => {
    const isTotal = row.category === 'Total'
    fTableRows.push(
      dataRow(
        [
          row.category,
          row.memberCount,
          row.monthlyPremium,
          row.annualPremium,
          row.totalAnnualPremium
        ],
        fColWidths,
        {
          alignments: fAlignments,
          bold: isTotal || undefined,
          fillColor: isTotal
            ? COLORS.lightFill
            : i % 2 === 1
              ? 'FCFDFE'
              : undefined
        }
      )
    )
  })

  children.push(
    new Table({
      width: { size: funeralCw, type: WidthType.DXA },
      columnWidths: fColWidths,
      rows: fTableRows
    })
  )

  return {
    properties: {
      page: {
        size: { ...A4, orientation: PageOrientation.LANDSCAPE },
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name) },
    footers: { default: makeFooter() },
    children
  }
}

/**
 * Section 5: Premium Breakdown per category.
 * Landscape orientation. One benefit table + funeral sub-table per category.
 */
function buildPremiumBreakdownSection(
  quote: any,
  resultSummaries: any[],
  benefitTitles: BenefitTitles
) {
  const hasNonFuneral = hasAnyNonFuneralBenefits(resultSummaries)
  const children: (Paragraph | Table)[] = [sectionHeading('Premium Breakdown')]

  const cw = LANDSCAPE_CONTENT_WIDTH

  resultSummaries.forEach((item) => {
    // Category heading
    if (hasNonFuneral) {
      children.push(categoryHeading(`${item.category} Category`))
    }

    // Benefit breakdown table (if category has non-funeral benefits)
    if (categoryHasNonFuneralBenefits(item)) {
      const breakdownRows = buildPremiumBreakdownRows(item, benefitTitles)
      const bCols = [
        'Benefit',
        'Total Sum Assured',
        'Annual Premium',
        '% Salary'
      ]
      const bColWidths = [
        Math.round(cw * 0.3),
        Math.round(cw * 0.25),
        Math.round(cw * 0.25),
        Math.round(cw * 0.2)
      ]
      bColWidths[3] = cw - bColWidths.slice(0, 3).reduce((a, b) => a + b, 0)

      const bAlignments = [
        AlignmentType.LEFT,
        AlignmentType.RIGHT,
        AlignmentType.RIGHT,
        AlignmentType.CENTER
      ]

      const rows = [headerRow(bCols, bColWidths, cw)]
      breakdownRows.forEach((row, i) => {
        rows.push(
          dataRow(
            [
              row.benefit,
              row.totalSumAssured,
              row.annualPremium,
              row.percentSalary
            ],
            bColWidths,
            {
              alignments: bAlignments,
              fillColor: i % 2 === 1 ? 'FDFDFE' : undefined
            }
          )
        )
      })

      children.push(
        new Table({
          width: { size: cw, type: WidthType.DXA },
          columnWidths: bColWidths,
          rows
        })
      )
    }

    // Group Funeral sub-table for this category
    children.push(
      new Paragraph({
        spacing: { before: 120, after: 40 },
        children: [
          new TextRun({
            text: `${item.category} - Group Funeral`,
            font: FONT,
            size: SIZE.body,
            bold: true,
            color: COLORS.dark
          })
        ]
      })
    )

    const funeralKV = buildGroupFuneralBreakdownRows(item)
    children.push(keyValueTable(funeralKV, Math.round(cw * 0.5)))

    // Spacer between categories
    children.push(new Paragraph({ spacing: { before: 100 }, children: [] }))
  })

  return {
    properties: {
      page: {
        size: { ...A4, orientation: PageOrientation.LANDSCAPE },
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name) },
    footers: { default: makeFooter() },
    children
  }
}

/**
 * Section 6: Benefits and Definitions of the Cover.
 * Landscape orientation. Per scheme_category: common benefits,
 * 7-column detail table, funeral coverage, optional educator benefits.
 */
function buildBenefitsDefinitionsSection(
  quote: any,
  resultSummaries: any[],
  categoryEducatorBenefits: any[],
  benefitTitles: BenefitTitles
) {
  const hasNonFuneral = hasAnyNonFuneralBenefits(resultSummaries)
  const children: (Paragraph | Table)[] = [
    sectionHeading('Benefits and Definitions of the Cover')
  ]

  const cw = LANDSCAPE_CONTENT_WIDTH

  const categories = quote.scheme_categories || []

  categories.forEach((item: any) => {
    children.push(categoryHeading(item.scheme_category))

    // Common benefits key-value table
    if (hasNonFuneral) {
      const commonRows = buildCategoryCommonBenefitRows(item, quote)
      children.push(keyValueTable(commonRows, Math.round(cw * 0.55)))
      children.push(new Paragraph({ spacing: { before: 60 }, children: [] }))
    }

    // 7-column benefit definitions table
    if (hasNonFuneral) {
      const defRows = buildBenefitDefinitionRows(item, quote, benefitTitles)
      const defCols = [
        'Benefit',
        'Salary Multiple',
        'Benefit Structure',
        'Waiting Period',
        'Deferred Period',
        'Cover Definition',
        'Risk Type'
      ]
      // Use smaller proportions for 7 columns
      const dColWidths = [
        Math.round(cw * 0.16),
        Math.round(cw * 0.12),
        Math.round(cw * 0.13),
        Math.round(cw * 0.12),
        Math.round(cw * 0.12),
        Math.round(cw * 0.18),
        Math.round(cw * 0.17)
      ]
      dColWidths[6] = cw - dColWidths.slice(0, 6).reduce((a, b) => a + b, 0)

      const defAlignments = [
        AlignmentType.LEFT,
        AlignmentType.CENTER,
        AlignmentType.CENTER,
        AlignmentType.CENTER,
        AlignmentType.CENTER,
        AlignmentType.CENTER,
        AlignmentType.CENTER
      ]

      const dRows = [headerRow(defCols, dColWidths, cw)]
      defRows.forEach((row, i) => {
        dRows.push(
          dataRow(
            [
              row.benefit,
              row.salaryMultiple,
              row.benefitStructure,
              row.waitingPeriod,
              row.deferredPeriod,
              row.coverDefinition,
              row.riskType
            ],
            dColWidths,
            {
              alignments: defAlignments,
              fillColor: i % 2 === 1 ? 'FCFDFE' : undefined
            }
          )
        )
      })

      children.push(
        new Table({
          width: { size: cw, type: WidthType.DXA },
          columnWidths: dColWidths,
          rows: dRows
        })
      )
      children.push(new Paragraph({ spacing: { before: 80 }, children: [] }))
    }

    // Group Funeral coverage table
    children.push(
      new Paragraph({
        spacing: { before: 80, after: 40 },
        children: [
          new TextRun({
            text: `${item.scheme_category} - Group Funeral`,
            font: FONT,
            size: SIZE.body,
            bold: true,
            color: COLORS.dark
          })
        ]
      })
    )

    const funeralCovRows = buildFuneralCoverageRows(item)
    const fCols = ['Member', 'Sum Assured', 'Maximum Number Covered']
    const fcw = Math.round(cw * 0.55)
    const fColWidths = [
      Math.round(fcw * 0.4),
      Math.round(fcw * 0.3),
      Math.round(fcw * 0.3)
    ]
    fColWidths[2] = fcw - fColWidths[0] - fColWidths[1]

    const fRows = [headerRow(fCols, fColWidths, fcw)]
    funeralCovRows.forEach((row, i) => {
      fRows.push(
        dataRow(
          [`${row.member}`, `${row.sumAssured}`, `${row.maxCovered}`],
          fColWidths,
          {
            alignments: [
              AlignmentType.LEFT,
              AlignmentType.RIGHT,
              AlignmentType.CENTER
            ],
            fillColor: i % 2 === 1 ? 'FCFDFE' : undefined
          }
        )
      )
    })

    children.push(
      new Table({
        width: { size: fcw, type: WidthType.DXA },
        columnWidths: fColWidths,
        rows: fRows
      })
    )

    // Educator Benefits (if present)
    if (hasNonFuneral && categoryEducatorBenefits.length > 0) {
      const eduRows = buildEducatorBenefitRows(
        categoryEducatorBenefits,
        item.scheme_category
      )
      if (eduRows.length > 0) {
        children.push(
          new Paragraph({
            spacing: { before: 120, after: 40 },
            children: [
              new TextRun({
                text: 'Educator Benefits',
                font: FONT,
                size: SIZE.body,
                bold: true,
                color: COLORS.dark
              })
            ]
          })
        )

        const eCols = [
          'Education Level',
          'Maximum Tuition per Year',
          'Maximum Coverage Period'
        ]
        const ecw = Math.round(cw * 0.55)
        const eColWidths = [
          Math.round(ecw * 0.4),
          Math.round(ecw * 0.35),
          Math.round(ecw * 0.25)
        ]
        eColWidths[2] = ecw - eColWidths[0] - eColWidths[1]

        const eRows = [headerRow(eCols, eColWidths, ecw)]
        eduRows.forEach((row, i) => {
          eRows.push(
            dataRow(
              [`${row.level}`, `${row.maxTuition}`, `${row.maxCoverage}`],
              eColWidths,
              {
                alignments: [
                  AlignmentType.LEFT,
                  AlignmentType.RIGHT,
                  AlignmentType.CENTER
                ],
                fillColor: i % 2 === 1 ? 'FCFDFE' : undefined
              }
            )
          )
        })

        children.push(
          new Table({
            width: { size: ecw, type: WidthType.DXA },
            columnWidths: eColWidths,
            rows: eRows
          })
        )
      }
    }

    // Spacer between categories
    children.push(new Paragraph({ spacing: { before: 200 }, children: [] }))
  })

  return {
    properties: {
      page: {
        size: { ...A4, orientation: PageOrientation.LANDSCAPE },
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name) },
    footers: { default: makeFooter() },
    children
  }
}

/**
 * Section 7: Underwriting and General Provisions.
 * Portrait orientation. Provisions text + closing block.
 */
function buildProvisionsSection(quote: any, insurer: any) {
  return {
    properties: {
      page: { size: A4, margin: MARGINS }
    },
    headers: { default: makeHeader(quote.scheme_name) },
    footers: { default: makeFooter() },
    children: [
      sectionHeading('Underwriting and General Provisions'),
      bodyText(insurer.general_provisions_text || ''),
      // Closing block
      new Paragraph({ spacing: { before: 400 }, children: [] }),
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.lightFill },
        spacing: { before: 0, after: 0 },
        children: [
          new TextRun({
            text: 'Thank you for considering our quotation.',
            font: FONT,
            size: SIZE.body,
            italics: true,
            color: COLORS.secondary
          })
        ]
      }),
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.lightFill },
        spacing: { before: 0, after: 0 },
        children: [
          new TextRun({
            text: 'We look forward to the opportunity to serve your insurance needs.',
            font: FONT,
            size: SIZE.body,
            italics: true,
            color: COLORS.secondary
          })
        ]
      }),
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.lightFill },
        spacing: { before: 80, after: 0 },
        children: [
          new TextRun({
            text: insurer.name,
            font: FONT,
            size: SIZE.body,
            bold: true,
            color: COLORS.secondary
          })
        ]
      })
    ]
  }
}

/**
 * Section 8: Acceptance of Quotation form.
 * Portrait orientation. Redesigned for Word using tables and form fields
 * instead of coordinate-based drawing from the PDF version.
 */
function buildAcceptanceFormSection(quote: any) {
  const cw = CONTENT_WIDTH

  // Helper: a form field row (label + underline input area)
  function formFieldRow(
    fields: { label: string; placeholder?: string }[]
  ): TableRow {
    const fieldWidth = Math.round(cw / fields.length)
    return new TableRow({
      children: fields.map(
        (f, i) =>
          new TableCell({
            width: {
              size:
                i === fields.length - 1
                  ? cw - fieldWidth * (fields.length - 1)
                  : fieldWidth,
              type: WidthType.DXA
            },
            borders: noBorders,
            margins: { top: 40, bottom: 40, left: 80, right: 80 },
            children: [
              new Paragraph({
                spacing: { after: 20 },
                children: [
                  new TextRun({
                    text: f.label,
                    font: FONT,
                    size: 12,
                    color: COLORS.mediumGray
                  })
                ]
              }),
              new Paragraph({
                border: {
                  bottom: {
                    style: BorderStyle.SINGLE,
                    size: 3,
                    color: COLORS.inputLine,
                    space: 1
                  }
                },
                spacing: { before: 80, after: 40 },
                children: f.placeholder
                  ? [
                      new TextRun({
                        text: f.placeholder,
                        font: FONT,
                        size: 12,
                        color: 'B4B4B4'
                      })
                    ]
                  : [new TextRun({ text: ' ', font: FONT, size: SIZE.body })]
              })
            ]
          })
      )
    })
  }

  // Helper: signature box
  function signatureBox(title: string): Table {
    const boxWidth = cw - 120
    return new Table({
      width: { size: boxWidth, type: WidthType.DXA },
      columnWidths: [boxWidth],
      rows: [
        new TableRow({
          children: [
            new TableCell({
              width: { size: boxWidth, type: WidthType.DXA },
              borders: thinBorders,
              margins: { top: 60, bottom: 60, left: 100, right: 100 },
              shading: { fill: COLORS.white, type: ShadingType.CLEAR },
              children: [
                new Paragraph({
                  spacing: { after: 20 },
                  children: [
                    new TextRun({
                      text: title,
                      font: FONT,
                      size: 14,
                      bold: true,
                      color: COLORS.navy
                    })
                  ]
                }),
                // Space for signature
                new Paragraph({
                  spacing: { before: 300, after: 0 },
                  children: []
                }),
                new Paragraph({
                  border: {
                    bottom: {
                      style: BorderStyle.SINGLE,
                      size: 3,
                      color: COLORS.inputLine,
                      space: 1
                    }
                  },
                  children: [
                    new TextRun({
                      text: 'Sign here',
                      font: FONT,
                      size: 12,
                      color: '969696'
                    })
                  ]
                })
              ]
            })
          ]
        })
      ]
    })
  }

  // Helper: section header with background
  function formSectionHeader(
    text: string,
    fillColor: string,
    textColor: string
  ): Paragraph {
    return new Paragraph({
      shading: { type: ShadingType.CLEAR, fill: fillColor },
      spacing: { before: 200, after: 80 },
      children: [
        new TextRun({
          text,
          font: FONT,
          size: 14,
          bold: true,
          color: textColor
        })
      ]
    })
  }

  return {
    properties: {
      page: { size: A4, margin: MARGINS }
    },
    headers: { default: makeHeader(quote.scheme_name) },
    footers: { default: makeFooter() },
    children: [
      // ===== HEADER BAR =====
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.navy },
        spacing: { before: 0, after: 0 },
        children: [
          new TextRun({
            text: 'ACCEPTANCE OF QUOTATION',
            font: FONT,
            size: SIZE.title,
            bold: true,
            color: COLORS.white
          })
        ]
      }),
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.navy },
        spacing: { before: 0, after: 100 },
        children: [
          new TextRun({
            text: 'POPIA Compliant',
            font: FONT,
            size: SIZE.small,
            color: 'E6E6E6'
          })
        ]
      }),

      // ===== POLICY DETAILS =====
      formSectionHeader('POLICY DETAILS', 'FAFBFC', COLORS.navy),
      new Table({
        width: { size: cw, type: WidthType.DXA },
        columnWidths: [Math.round(cw / 2), cw - Math.round(cw / 2)],
        rows: [
          formFieldRow([
            { label: 'EMPLOYER / SCHEME NAME' },
            { label: 'QUOTE NUMBER' }
          ]),
          formFieldRow([
            { label: 'DATE OF QUOTE', placeholder: 'DD/MM/YYYY' },
            { label: 'COMMENCEMENT DATE', placeholder: 'DD/MM/YYYY' }
          ])
        ]
      }),

      // ===== PROFILE NOTICE =====
      new Paragraph({ spacing: { before: 160 }, children: [] }),
      bodyText(
        'If the member data profile at the quotation implementation date differ by 7% or more from that on which the quotation was based, we reserve the right to revise the rates and Automatic Acceptance Limit. The Employer/Scheme will be notified accordingly and must provide acceptance before implementation proceeds.'
      ),
      bodyText(
        'By signing this quotation, the Employer/Scheme acknowledges that they have read, understood, and agree to be bound by all the terms and conditions of this quotation.'
      ),

      // ===== EMPLOYER AUTHORISATION =====
      formSectionHeader('EMPLOYER AUTHORISATION', 'FAFBFC', COLORS.navy),
      signatureBox('Duly Authorised Signatory'),
      new Paragraph({ spacing: { before: 60 }, children: [] }),
      new Table({
        width: { size: cw, type: WidthType.DXA },
        columnWidths: [
          Math.round(cw * 0.4),
          Math.round(cw * 0.3),
          cw - Math.round(cw * 0.4) - Math.round(cw * 0.3)
        ],
        rows: [
          formFieldRow([
            { label: 'FULL NAME' },
            { label: 'CAPACITY', placeholder: 'e.g., Director' },
            { label: 'DATE', placeholder: 'DD/MM/YYYY' }
          ])
        ]
      }),

      // ===== INTERMEDIARY DETAILS =====
      formSectionHeader('INTERMEDIARY DETAILS', 'FAFBFC', COLORS.navy),
      signatureBox('Intermediary / FAIS Representative'),
      new Paragraph({ spacing: { before: 60 }, children: [] }),
      new Table({
        width: { size: cw, type: WidthType.DXA },
        columnWidths: [
          Math.round(cw * 0.4),
          Math.round(cw * 0.3),
          cw - Math.round(cw * 0.4) - Math.round(cw * 0.3)
        ],
        rows: [
          formFieldRow([
            { label: 'FULL NAME' },
            { label: 'FAIS REG NO.' },
            { label: 'DATE', placeholder: 'DD/MM/YYYY' }
          ])
        ]
      }),

      // ===== POPIA CONSENT =====
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.lightBlue },
        spacing: { before: 200, after: 40 },
        children: [
          new TextRun({
            text: 'POPIA Consent & Data Protection',
            font: FONT,
            size: 14,
            bold: true,
            color: '0366D6'
          })
        ]
      }),
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.lightBlue },
        spacing: { after: 60 },
        children: [
          new TextRun({
            text: 'In terms of the Protection of Personal Information Act 4 of 2013 (POPIA), the Employer consents to the processing of personal information of employees and scheme members for the purpose of underwriting, administering, and processing claims under this Group Risk policy. Information will be processed lawfully, minimally, and only for the specific purpose stated.',
            font: FONT,
            size: 14,
            color: COLORS.dark
          })
        ]
      }),
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.lightBlue },
        spacing: { after: 80 },
        children: [
          new TextRun({
            text: '\u2610  I confirm that the Employer has obtained necessary consent from data subjects (employees/members) for the processing of their personal information as required by POPIA, and warrants that all information provided is true and complete.',
            font: FONT,
            size: 13,
            color: COLORS.dark
          })
        ]
      }),

      // ===== FOR OFFICE USE ONLY =====
      new Paragraph({
        shading: { type: ShadingType.CLEAR, fill: COLORS.lightOrange },
        spacing: { before: 200, after: 80 },
        children: [
          new TextRun({
            text: 'FOR OFFICE USE ONLY',
            font: FONT,
            size: 14,
            bold: true,
            color: COLORS.orange
          })
        ]
      }),
      (() => {
        const colW = Math.round(cw / 3)
        const lastColW = cw - colW * 2
        const officeBorder = {
          style: BorderStyle.SINGLE,
          size: 3,
          color: 'FFB47C'
        }
        const officeBorders = {
          top: officeBorder,
          bottom: officeBorder,
          left: officeBorder,
          right: officeBorder
        }
        return new Table({
          width: { size: cw, type: WidthType.DXA },
          columnWidths: [colW, colW, lastColW],
          rows: [
            new TableRow({
              children: ['RECEIVED BY', 'DATE RECEIVED', 'POLICY NUMBER'].map(
                (label, i) =>
                  new TableCell({
                    width: {
                      size: i === 2 ? lastColW : colW,
                      type: WidthType.DXA
                    },
                    borders: officeBorders,
                    shading: {
                      fill: COLORS.lightOrange,
                      type: ShadingType.CLEAR
                    },
                    margins: cellMargins,
                    children: [
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: label,
                            font: FONT,
                            size: 12,
                            color: COLORS.mediumGray
                          })
                        ]
                      })
                    ]
                  })
              )
            }),
            new TableRow({
              children: ['UNDERWRITER', 'APPROVED BY', 'DATE APPROVED'].map(
                (label, i) =>
                  new TableCell({
                    width: {
                      size: i === 2 ? lastColW : colW,
                      type: WidthType.DXA
                    },
                    borders: officeBorders,
                    shading: {
                      fill: COLORS.lightOrange,
                      type: ShadingType.CLEAR
                    },
                    margins: cellMargins,
                    children: [
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: label,
                            font: FONT,
                            size: 12,
                            color: COLORS.mediumGray
                          })
                        ]
                      })
                    ]
                  })
              )
            })
          ]
        })
      })()
    ]
  }
}

// ---------------------------------------------------------------------------
// Main composable
// ---------------------------------------------------------------------------

export function useDocxQuoteGeneration() {
  const isGenerating = ref(false)
  const errorMessage = ref('')

  async function generateDocxQuote(
    quote: any,
    resultSummaries: any[],
    insurer: any,
    categoryEducatorBenefits: any[],
    benefitMaps: any[],
    benefitTitles: BenefitTitles
  ): Promise<void> {
    isGenerating.value = true
    errorMessage.value = ''
    try {
      const doc = new Document({
        styles: {
          default: {
            document: {
              run: { font: FONT, size: SIZE.body }
            }
          }
        },
        sections: [
          buildCoverAndSummarySection(quote, resultSummaries, insurer),
          buildPremiumSummarySection(quote, resultSummaries),
          buildPremiumBreakdownSection(quote, resultSummaries, benefitTitles),
          buildBenefitsDefinitionsSection(
            quote,
            resultSummaries,
            categoryEducatorBenefits,
            benefitTitles
          ),
          buildProvisionsSection(quote, insurer),
          buildAcceptanceFormSection(quote)
        ]
      })

      const blob = await Packer.toBlob(doc)
      const filename = `${quote.scheme_name}_Quotation_${formatDateString(new Date(), true, true, true)}.docx`
      saveAs(blob, filename)
    } catch (err: any) {
      console.error('DOCX generation error:', err)
      errorMessage.value = `Error generating Word document: ${err.message}`
    } finally {
      isGenerating.value = false
    }
  }

  return {
    isGenerating,
    errorMessage,
    generateDocxQuote
  }
}
