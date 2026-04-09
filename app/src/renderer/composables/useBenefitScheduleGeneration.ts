/**
 * Composable for generating a Benefit Schedule document.
 *
 * Consolidates benefit configuration (salary multiples, waiting/deferred periods,
 * cover definitions) and premium summary (sum assured, annual premium, % of salary)
 * into a single downloadable PDF or Word document.
 *
 * Reuses design constants and helper functions from useDocxQuoteGeneration.ts
 * and quoteDataHelpers.ts.
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
  AlignmentType,
  BorderStyle,
  WidthType,
  ShadingType,
  Header,
  Footer,
  PageNumber,
  TabStopType,
  TabStopPosition
} from 'docx'
import { saveAs } from 'file-saver'
import jsPDF from 'jspdf'
import { applyPlugin } from 'jspdf-autotable'
import formatDateString from '@/renderer/utils/helpers.js'
import type { BenefitTitles } from '@/renderer/types/docxQuote'
import {
  buildBenefitDefinitionRows,
  buildPremiumBreakdownRows,
  buildFuneralCoverageRows,
  roundUpToTwoDecimalsAccounting
} from '@/renderer/utils/quoteDataHelpers'

applyPlugin(jsPDF)

// ---------------------------------------------------------------------------
// Design constants (matching useDocxQuoteGeneration.ts)
// ---------------------------------------------------------------------------

const COLORS = {
  primary: '34495E',
  secondary: '34495E',
  lightFill: 'ECF0F1',
  dark: '2C3E50',
  altRow: 'FAFAFA',
  white: 'FFFFFF'
}

const FONT = 'Arial'
const A4 = { width: 11906, height: 16838 }
const MARGINS = { top: 850, bottom: 850, left: 1020, right: 1020 }
const CONTENT_WIDTH = A4.width - MARGINS.left - MARGINS.right
const LANDSCAPE_CONTENT_WIDTH = A4.height - MARGINS.left - MARGINS.right

const SIZE = {
  heading: 28,
  subheading: 24,
  body: 20,
  caption: 18
}

const thinBorder = { style: BorderStyle.SINGLE, size: 4, color: 'D5D8DC' }
const thinBorders = {
  top: thinBorder,
  bottom: thinBorder,
  left: thinBorder,
  right: thinBorder
}
const cellMargins = { top: 60, bottom: 60, left: 100, right: 100 }

// ---------------------------------------------------------------------------
// Shared DOCX element builders
// ---------------------------------------------------------------------------

function makeHeader(schemeName: string): Header {
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
            text: `${schemeName} – Benefit Schedule`,
            font: FONT,
            size: SIZE.caption,
            color: COLORS.secondary
          })
        ]
      })
    ]
  })
}

function makeFooter(): Footer {
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

function spacer(): Paragraph {
  return new Paragraph({ spacing: { after: 80 }, children: [] })
}

function tableHeaderRow(labels: string[], colWidths: number[]): TableRow {
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

function tableDataRow(
  values: string[],
  colWidths: number[],
  options?: { bold?: boolean; fillColor?: string }
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
              alignment: i === 0 ? AlignmentType.LEFT : AlignmentType.RIGHT,
              children: [
                new TextRun({
                  text: val ?? '-',
                  font: FONT,
                  size: SIZE.body,
                  bold: options?.bold ?? false,
                  color: COLORS.dark
                })
              ]
            })
          ]
        })
    )
  })
}

// ---------------------------------------------------------------------------
// DOCX section builders
// ---------------------------------------------------------------------------

function buildSchemeInfoSection(quote: any, resultSummaries: any[]): any {
  const totalLives = resultSummaries.reduce(
    (s, r) => s + (r.member_count || 0),
    0
  )
  const totalPremium = resultSummaries.reduce(
    (s, r) => s + (r.exp_total_annual_premium_excl_funeral || 0),
    0
  )
  const categories = (quote.selected_scheme_categories || []).join(', ') || '-'

  const kvColWidths = [
    Math.round(CONTENT_WIDTH * 0.35),
    Math.round(CONTENT_WIDTH * 0.65)
  ]

  const infoRows = [
    ['Scheme Name', quote.scheme_name || '-'],
    ['Quote Reference', `${quote.id || '-'}`],
    [
      'Commencement Date',
      formatDateString(quote.effective_date, true, false, false) || '-'
    ],
    [
      'Prepared Date',
      formatDateString(quote.creation_date, true, false, false) || '-'
    ],
    ['Total Lives', `${totalLives}`],
    [
      'Total Annual Premium (excl. Funeral)',
      roundUpToTwoDecimalsAccounting(totalPremium)
    ],
    ['Categories', categories]
  ]

  return {
    properties: {
      page: {
        size: { width: A4.width, height: A4.height },
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name || 'Benefit Schedule') },
    footers: { default: makeFooter() },
    children: [
      sectionHeading('Benefit Schedule'),
      new Table({
        width: { size: CONTENT_WIDTH, type: WidthType.DXA },
        rows: infoRows.map(
          ([label, value]) =>
            new TableRow({
              children: [
                new TableCell({
                  width: { size: kvColWidths[0], type: WidthType.DXA },
                  shading: { fill: COLORS.lightFill, type: ShadingType.CLEAR },
                  borders: thinBorders,
                  margins: cellMargins,
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: label,
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
                  width: { size: kvColWidths[1], type: WidthType.DXA },
                  borders: thinBorders,
                  margins: cellMargins,
                  children: [
                    new Paragraph({
                      children: [
                        new TextRun({
                          text: value,
                          font: FONT,
                          size: SIZE.body,
                          color: COLORS.dark
                        })
                      ]
                    })
                  ]
                })
              ]
            })
        )
      })
    ]
  }
}

function buildBenefitConfigSection(
  quote: any,
  resultSummaries: any[],
  benefitTitles: BenefitTitles
): any {
  const colWidths = [
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.15), // Category
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.12), // Benefit
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.11), // Salary Multiple
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.12), // Benefit Structure
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.12), // Waiting Period
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.12), // Deferred Period
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.14), // Cover Definition
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.12) // Risk Type
  ]

  const headers = [
    'Category',
    'Benefit',
    'Salary Multiple',
    'Benefit Structure',
    'Waiting Period',
    'Deferred Period',
    'Cover Definition',
    'Risk Type'
  ]

  const separatorRow = (widths: number[]): TableRow =>
    new TableRow({
      height: { value: 120, rule: 'exact' as any },
      children: widths.map(
        (w) =>
          new TableCell({
            width: { size: w, type: WidthType.DXA },
            shading: { fill: COLORS.primary, type: ShadingType.CLEAR },
            borders: {
              top: { style: BorderStyle.NONE, size: 0, color: 'FFFFFF' },
              bottom: { style: BorderStyle.NONE, size: 0, color: 'FFFFFF' },
              left: { style: BorderStyle.NONE, size: 0, color: 'FFFFFF' },
              right: { style: BorderStyle.NONE, size: 0, color: 'FFFFFF' }
            },
            margins: { top: 0, bottom: 0, left: 0, right: 0 },
            children: [new Paragraph({ children: [] })]
          })
      )
    })

  const rows: TableRow[] = [tableHeaderRow(headers, colWidths)]

  resultSummaries.forEach((item, catIdx) => {
    if (catIdx > 0) {
      rows.push(separatorRow(colWidths))
    }
    const defRows = buildBenefitDefinitionRows(item, quote, benefitTitles)
    defRows.forEach((r, idx) => {
      rows.push(
        tableDataRow(
          [
            idx === 0 ? item.category || '-' : '',
            r.benefit,
            r.salaryMultiple,
            r.benefitStructure,
            r.waitingPeriod,
            r.deferredPeriod,
            r.coverDefinition,
            r.riskType
          ],
          colWidths,
          { fillColor: idx % 2 === 0 ? COLORS.white : COLORS.altRow }
        )
      )
    })
  })

  return {
    properties: {
      page: {
        size: { width: A4.height, height: A4.width }, // landscape
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name || 'Benefit Schedule') },
    footers: { default: makeFooter() },
    children: [
      sectionHeading('Benefit Configuration'),
      spacer(),
      new Table({
        width: { size: LANDSCAPE_CONTENT_WIDTH, type: WidthType.DXA },
        rows
      })
    ]
  }
}

function buildPremiumSummarySection(
  quote: any,
  resultSummaries: any[],
  benefitTitles: BenefitTitles
): any {
  const colWidths = [
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.18),
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.18),
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.22),
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.22),
    Math.round(LANDSCAPE_CONTENT_WIDTH * 0.2)
  ]
  const headers = [
    'Category',
    'Benefit',
    'Total Sum Assured',
    'Annual Premium',
    '% of Salary'
  ]

  const rows: TableRow[] = [tableHeaderRow(headers, colWidths)]

  resultSummaries.forEach((item) => {
    const breakdownRows = buildPremiumBreakdownRows(item, benefitTitles)
    breakdownRows.forEach((r, idx) => {
      rows.push(
        tableDataRow(
          [
            idx === 0 ? item.category || '-' : '',
            r.benefit,
            r.totalSumAssured,
            r.annualPremium,
            r.percentSalary
          ],
          colWidths,
          { fillColor: idx % 2 === 0 ? COLORS.white : COLORS.altRow }
        )
      )
    })

    // Sub-total row
    rows.push(
      tableDataRow(
        [
          '',
          'Sub Total (excl. Funeral)',
          '-',
          roundUpToTwoDecimalsAccounting(
            item.exp_total_annual_premium_excl_funeral || 0
          ),
          `${roundUpToTwoDecimalsAccounting((item.proportion_exp_total_premium_excl_funeral_salary || 0) * 100)}%`
        ],
        colWidths,
        { bold: true, fillColor: COLORS.lightFill }
      )
    )
  })

  const funeralSections: Paragraph[] = []
  const hasFuneral = resultSummaries.some(
    (r) => (r.family_funeral_main_member_funeral_sum_assured || 0) > 0
  )
  if (hasFuneral) {
    funeralSections.push(spacer(), sectionHeading('Funeral Coverage'))
    resultSummaries.forEach((item) => {
      const fRows = buildFuneralCoverageRows(item)
      const fColWidths = [
        Math.round(LANDSCAPE_CONTENT_WIDTH * 0.3),
        Math.round(LANDSCAPE_CONTENT_WIDTH * 0.35),
        Math.round(LANDSCAPE_CONTENT_WIDTH * 0.35)
      ]
      funeralSections.push(categoryHeading(item.category || ''), spacer())
      const fTableRows: TableRow[] = [
        tableHeaderRow(
          ['Member Type', 'Sum Assured', 'Max Covered'],
          fColWidths
        ),
        ...fRows.map((r) =>
          tableDataRow(
            [r.member, `${r.sumAssured}`, `${r.maxCovered}`],
            fColWidths
          )
        )
      ]
      funeralSections.push(
        new Table({
          width: { size: LANDSCAPE_CONTENT_WIDTH, type: WidthType.DXA },
          rows: fTableRows
        }) as unknown as Paragraph
      )
    })
  }

  return {
    properties: {
      page: {
        size: { width: A4.height, height: A4.width }, // landscape
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name || 'Benefit Schedule') },
    footers: { default: makeFooter() },
    children: [
      sectionHeading('Premium Summary'),
      spacer(),
      new Table({
        width: { size: LANDSCAPE_CONTENT_WIDTH, type: WidthType.DXA },
        rows
      }),
      ...funeralSections
    ]
  }
}

// ---------------------------------------------------------------------------
// PDF generation helpers
// ---------------------------------------------------------------------------

function pdfColors() {
  return {
    primary: [52, 73, 94] as [number, number, number],
    light: [236, 240, 241] as [number, number, number],
    dark: [44, 62, 80] as [number, number, number],
    white: [255, 255, 255] as [number, number, number],
    altRow: [250, 250, 250] as [number, number, number]
  }
}

function addPdfHeaderFooter(doc: any, schemeName: string, pageNumber: number) {
  const c = pdfColors()
  const pw = doc.internal.pageSize.getWidth()
  const ph = doc.internal.pageSize.getHeight()
  doc.setFillColor(...c.light)
  doc.rect(0, 0, pw, 15, 'F')
  doc.setTextColor(...c.primary)
  doc.setFontSize(9)
  doc.setFont('helvetica', 'normal')
  doc.text(`${schemeName} – Benefit Schedule`, 14, 10)
  doc.setFillColor(...c.light)
  doc.rect(0, ph - 15, pw, 15, 'F')
  doc.text(
    `Generated on ${formatDateString(new Date(), true, true, true)}`,
    14,
    ph - 7
  )
  doc.text(`Page ${pageNumber}`, pw - 25, ph - 7)
}

// ---------------------------------------------------------------------------
// Public composable
// ---------------------------------------------------------------------------

export function useBenefitScheduleGeneration() {
  const isGenerating = ref(false)

  async function generateBenefitScheduleDocx(
    quote: any,
    resultSummaries: any[],
    benefitTitles: BenefitTitles
  ): Promise<void> {
    isGenerating.value = true
    try {
      const doc = new Document({
        sections: [
          buildSchemeInfoSection(quote, resultSummaries),
          buildBenefitConfigSection(quote, resultSummaries, benefitTitles),
          buildPremiumSummarySection(quote, resultSummaries, benefitTitles)
        ]
      })
      const blob = await Packer.toBlob(doc)
      const fileName =
        `${quote.scheme_name || 'Benefit_Schedule'}_Benefit_Schedule.docx`.replace(
          /\s+/g,
          '_'
        )
      saveAs(blob, fileName)
    } finally {
      isGenerating.value = false
    }
  }

  async function generateBenefitSchedulePdf(
    quote: any,
    resultSummaries: any[],
    benefitTitles: BenefitTitles
  ): Promise<void> {
    isGenerating.value = true
    try {
      // eslint-disable-next-line new-cap
      const doc: any = new jsPDF()
      const c = pdfColors()
      let pageNumber = 1
      const schemeName = quote.scheme_name || 'Benefit Schedule'
      const leftMargin = 14
      const topMargin = 20
      let currentY = topMargin

      // --- Page 1: Scheme Info ---
      addPdfHeaderFooter(doc, schemeName, pageNumber)
      currentY = topMargin + 5

      doc.setFontSize(14)
      doc.setFont('helvetica', 'bold')
      doc.setTextColor(...c.primary)
      doc.text('Benefit Schedule', leftMargin, currentY)
      currentY += 10

      const totalLives = resultSummaries.reduce(
        (s, r) => s + (r.member_count || 0),
        0
      )
      const totalPremium = resultSummaries.reduce(
        (s, r) => s + (r.exp_total_annual_premium_excl_funeral || 0),
        0
      )
      const categories =
        (quote.selected_scheme_categories || []).join(', ') || '-'

      const infoBody = [
        ['Scheme Name', quote.scheme_name || '-'],
        ['Quote Reference', `${quote.id || '-'}`],
        [
          'Commencement Date',
          formatDateString(quote.effective_date, true, false, false) || '-'
        ],
        [
          'Prepared Date',
          formatDateString(quote.creation_date, true, false, false) || '-'
        ],
        ['Total Lives', `${totalLives}`],
        [
          'Total Annual Premium (excl. Funeral)',
          roundUpToTwoDecimalsAccounting(totalPremium)
        ],
        ['Categories', categories]
      ]

      doc.autoTable({
        startY: currentY,
        body: infoBody,
        theme: 'grid',
        styles: {
          fontSize: 10,
          cellPadding: { top: 4, right: 8, bottom: 4, left: 8 }
        },
        columnStyles: {
          0: {
            fontStyle: 'bold',
            fillColor: c.light,
            textColor: c.dark,
            halign: 'left'
          },
          1: { textColor: c.dark }
        },
        didDrawPage: () => addPdfHeaderFooter(doc, schemeName, pageNumber)
      })

      // --- Page 2: Benefit Configuration (landscape) ---
      doc.addPage('a4', 'landscape')
      pageNumber++
      addPdfHeaderFooter(doc, schemeName, pageNumber)
      currentY = topMargin + 5

      doc.setFontSize(14)
      doc.setFont('helvetica', 'bold')
      doc.setTextColor(...c.primary)
      doc.text('Benefit Configuration', leftMargin, currentY)
      currentY += 8

      const configHead = [
        [
          'Category',
          'Benefit',
          'Salary Multiple',
          'Benefit Structure',
          'Waiting Period',
          'Deferred Period',
          'Cover Definition',
          'Risk Type'
        ]
      ]
      const configBody: any[][] = []

      resultSummaries.forEach((item, catIdx) => {
        if (catIdx > 0) {
          // narrow separator row between categories
          configBody.push([
            {
              content: '',
              colSpan: 8,
              styles: { fillColor: c.primary, minCellHeight: 3, cellPadding: 0 }
            }
          ])
        }
        const defRows = buildBenefitDefinitionRows(item, quote, benefitTitles)
        defRows.forEach((r, idx) => {
          configBody.push([
            idx === 0 ? item.category || '-' : '',
            r.benefit,
            r.salaryMultiple,
            r.benefitStructure,
            r.waitingPeriod,
            r.deferredPeriod,
            r.coverDefinition,
            r.riskType
          ])
        })
      })

      doc.autoTable({
        startY: currentY,
        head: configHead,
        body: configBody,
        theme: 'grid',
        styles: {
          fontSize: 8,
          cellPadding: { top: 3, right: 4, bottom: 3, left: 4 }
        },
        headStyles: {
          fillColor: c.primary,
          textColor: c.white,
          fontStyle: 'bold'
        },
        alternateRowStyles: { fillColor: c.altRow },
        didDrawPage: () => {
          pageNumber++
          addPdfHeaderFooter(doc, schemeName, pageNumber)
        }
      })

      // --- Page 3: Premium Summary (landscape) ---
      doc.addPage('a4', 'landscape')
      pageNumber++
      addPdfHeaderFooter(doc, schemeName, pageNumber)
      currentY = topMargin + 5

      doc.setFontSize(14)
      doc.setFont('helvetica', 'bold')
      doc.setTextColor(...c.primary)
      doc.text('Premium Summary', leftMargin, currentY)
      currentY += 8

      const premHead = [
        [
          'Category',
          'Benefit',
          'Total Sum Assured',
          'Annual Premium',
          '% of Salary'
        ]
      ]
      const premBody: any[][] = []

      resultSummaries.forEach((item) => {
        const bRows = buildPremiumBreakdownRows(item, benefitTitles)
        bRows.forEach((r, idx) => {
          premBody.push([
            idx === 0 ? item.category || '-' : '',
            r.benefit,
            { content: r.totalSumAssured, styles: { halign: 'right' } },
            { content: r.annualPremium, styles: { halign: 'right' } },
            { content: r.percentSalary, styles: { halign: 'right' } }
          ])
        })
        premBody.push([
          '',
          {
            content: 'Sub Total (excl. Funeral)',
            styles: { fontStyle: 'bold' }
          },
          '',
          {
            content: roundUpToTwoDecimalsAccounting(
              item.exp_total_annual_premium_excl_funeral || 0
            ),
            styles: { halign: 'right', fontStyle: 'bold' }
          },
          {
            content: `${roundUpToTwoDecimalsAccounting((item.proportion_exp_total_premium_excl_funeral_salary || 0) * 100)}%`,
            styles: { halign: 'right', fontStyle: 'bold' }
          }
        ])
      })

      doc.autoTable({
        startY: currentY,
        head: premHead,
        body: premBody,
        theme: 'grid',
        styles: {
          fontSize: 9,
          cellPadding: { top: 3, right: 6, bottom: 3, left: 6 }
        },
        headStyles: {
          fillColor: c.primary,
          textColor: c.white,
          fontStyle: 'bold'
        },
        alternateRowStyles: { fillColor: c.altRow },
        didDrawPage: () => {
          pageNumber++
          addPdfHeaderFooter(doc, schemeName, pageNumber)
        }
      })

      // --- Funeral Coverage (same page, if applicable) ---
      const hasFuneral = resultSummaries.some(
        (r) => (r.family_funeral_main_member_funeral_sum_assured || 0) > 0
      )
      if (hasFuneral) {
        const afterY = (doc as any).lastAutoTable.finalY + 12
        doc.setFontSize(12)
        doc.setFont('helvetica', 'bold')
        doc.setTextColor(...c.primary)
        doc.text('Funeral Coverage', leftMargin, afterY)

        const funeralHead = [
          ['Category', 'Member Type', 'Sum Assured', 'Max Covered']
        ]
        const funeralBody: string[][] = []

        resultSummaries.forEach((item) => {
          const fRows = buildFuneralCoverageRows(item)
          fRows.forEach((r, idx) => {
            funeralBody.push([
              idx === 0 ? item.category || '-' : '',
              r.member,
              `${r.sumAssured}`,
              `${r.maxCovered}`
            ])
          })
        })

        doc.autoTable({
          startY: afterY + 4,
          head: funeralHead,
          body: funeralBody,
          theme: 'grid',
          styles: {
            fontSize: 9,
            cellPadding: { top: 3, right: 6, bottom: 3, left: 6 }
          },
          headStyles: {
            fillColor: c.primary,
            textColor: c.white,
            fontStyle: 'bold'
          },
          alternateRowStyles: { fillColor: c.altRow },
          didDrawPage: () => {
            pageNumber++
            addPdfHeaderFooter(doc, schemeName, pageNumber)
          }
        })
      }

      const fileName = `${schemeName}_Benefit_Schedule.pdf`.replace(/\s+/g, '_')
      doc.save(fileName)
    } finally {
      isGenerating.value = false
    }
  }

  return {
    isGenerating,
    generateBenefitScheduleDocx,
    generateBenefitSchedulePdf
  }
}
