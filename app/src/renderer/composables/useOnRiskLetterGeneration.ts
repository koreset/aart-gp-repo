/**
 * Composable for generating On Risk Letter documents in Word (.docx) and PDF format.
 *
 * The On Risk letter is a formal confirmation issued to a scheme when a
 * group risk quote is accepted, confirming that cover is now active.
 *
 * Uses the npm `docx` package (same approach as useDocxQuoteGeneration.ts)
 * and `jspdf` + `jspdf-autotable` for PDF output.
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
  BorderStyle,
  WidthType,
  ShadingType,
  PageNumber,
  TabStopType,
  TabStopPosition
} from 'docx'
import { saveAs } from 'file-saver'
import jsPDF from 'jspdf'
import { applyPlugin } from 'jspdf-autotable'
import formatDateString from '@/renderer/utils/helpers.js'

applyPlugin(jsPDF)

// ---------------------------------------------------------------------------
// Design constants — matching the quotation document colour palette
// ---------------------------------------------------------------------------

const COLORS = {
  primary: '34495E',
  accent: 'E74C3C',
  lightFill: 'ECF0F1',
  dark: '2C3E50',
  white: 'FFFFFF',
  navy: '1E3A5F',
  mediumGray: '586069'
}

const FONT = 'Arial'

// A4 portrait in DXA
const A4 = { width: 11906, height: 16838 }
const MARGINS = { top: 1134, bottom: 1134, left: 1134, right: 1134 }
const CONTENT_WIDTH = A4.width - MARGINS.left - MARGINS.right

// Half-point sizes (multiply pt by 2)
const SIZE = {
  title: 32,
  heading: 28,
  subheading: 24,
  body: 20,
  caption: 18,
  small: 16
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
// Currency formatter
// ---------------------------------------------------------------------------

function formatCurrency(value: number, currency = 'ZAR'): string {
  const symbol = currency === 'ZAR' ? 'R' : currency
  return `${symbol} ${value.toLocaleString('en-ZA', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

// ---------------------------------------------------------------------------
// Header & Footer
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
            text: `${schemeName} — On Risk Letter`,
            font: FONT,
            size: SIZE.caption,
            color: COLORS.primary
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
            color: COLORS.mediumGray
          }),
          new TextRun({
            text: '\tPage ',
            font: FONT,
            size: SIZE.caption,
            color: COLORS.mediumGray
          }),
          new TextRun({
            children: [PageNumber.CURRENT],
            font: FONT,
            size: SIZE.caption,
            color: COLORS.mediumGray
          })
        ]
      })
    ]
  })
}

// ---------------------------------------------------------------------------
// Reusable element builders
// ---------------------------------------------------------------------------

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

function bodyText(
  text: string,
  options?: { bold?: boolean; italic?: boolean }
): Paragraph {
  return new Paragraph({
    alignment: AlignmentType.JUSTIFIED,
    spacing: { after: 100 },
    children: [
      new TextRun({
        text,
        font: FONT,
        size: SIZE.body,
        color: COLORS.dark,
        bold: options?.bold ?? false,
        italics: options?.italic ?? false
      })
    ]
  })
}

/** Two-column key/value row for the details table. */
function kvRow(label: string, value: string, colWidths: number[]): TableRow {
  return new TableRow({
    children: [
      new TableCell({
        width: { size: colWidths[0], type: WidthType.DXA },
        borders: thinBorders,
        margins: cellMargins,
        shading: { fill: COLORS.lightFill, type: ShadingType.CLEAR },
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
        width: { size: colWidths[1], type: WidthType.DXA },
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
}

// ---------------------------------------------------------------------------
// Document section builders
// ---------------------------------------------------------------------------

function buildLetterSection(data: any) {
  const {
    quote,
    scheme,
    insurer,
    letter,
    benefit_summary: benefitSummary
  } = data // eslint-disable-line camelcase

  const children: any[] = []

  // ── Insurer logo ──
  if (insurer.logo) {
    try {
      const logoFormat = (insurer.logo_mime_type || 'image/png')
        .split('/')[1]
        ?.toLowerCase()
      const logoBuffer = Uint8Array.from(atob(insurer.logo), (c) =>
        c.charCodeAt(0)
      )
      children.push(
        new Paragraph({
          spacing: { after: 200 },
          children: [
            new ImageRun({
              type: logoFormat === 'jpeg' ? 'jpg' : logoFormat,
              data: logoBuffer,
              transformation: { width: 180, height: 60 },
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
      console.warn('Failed to decode logo for On Risk letter:', e)
    }
  }

  // ── Insurer address block ──
  const addressLines = [
    insurer.name,
    insurer.address_line_1,
    insurer.address_line_2,
    insurer.address_line_3,
    [insurer.city, insurer.province, insurer.post_code]
      .filter(Boolean)
      .join(', '),
    insurer.telephone ? `Tel: ${insurer.telephone}` : '',
    insurer.email ? `Email: ${insurer.email}` : ''
  ].filter(Boolean)

  for (const line of addressLines) {
    children.push(
      new Paragraph({
        spacing: { after: 20 },
        children: [
          new TextRun({
            text: line,
            font: FONT,
            size: SIZE.body,
            color: COLORS.dark
          })
        ]
      })
    )
  }

  children.push(new Paragraph({ spacing: { after: 200 }, children: [] }))

  // ── Date ──
  children.push(
    new Paragraph({
      spacing: { after: 200 },
      children: [
        new TextRun({
          text: formatDateString(
            letter.letter_date || new Date(),
            true,
            true,
            true
          ),
          font: FONT,
          size: SIZE.body,
          color: COLORS.dark
        })
      ]
    })
  )

  // ── Addressee ──
  const addressee = [
    quote.scheme_name,
    quote.scheme_contact || scheme.contact_person || '',
    quote.scheme_email || scheme.contact_email || ''
  ].filter(Boolean)

  for (const line of addressee) {
    children.push(
      new Paragraph({
        spacing: { after: 20 },
        children: [
          new TextRun({
            text: line,
            font: FONT,
            size: SIZE.body,
            color: COLORS.dark
          })
        ]
      })
    )
  }

  // Broker details if applicable
  if (
    quote.distribution_channel === 'broker' &&
    (quote.quote_broker?.broker_name || quote.broker_name)
  ) {
    const brokerName =
      quote.quote_broker?.broker_name || quote.broker_name || ''
    if (brokerName) {
      children.push(
        new Paragraph({
          spacing: { after: 20 },
          children: [
            new TextRun({
              text: `Broker: ${brokerName}`,
              font: FONT,
              size: SIZE.body,
              color: COLORS.dark
            })
          ]
        })
      )
    }
  }

  children.push(new Paragraph({ spacing: { after: 200 }, children: [] }))

  // ── Subject line ──
  children.push(
    new Paragraph({
      spacing: { after: 200 },
      children: [
        new TextRun({
          text: 'Confirmation of Cover — On Risk Letter',
          font: FONT,
          size: SIZE.subheading,
          bold: true,
          color: COLORS.navy
        })
      ]
    })
  )

  if (letter.letter_reference) {
    children.push(
      new Paragraph({
        spacing: { after: 200 },
        children: [
          new TextRun({
            text: `Reference: ${letter.letter_reference}`,
            font: FONT,
            size: SIZE.body,
            color: COLORS.mediumGray
          })
        ]
      })
    )
  }

  // ── Opening paragraph ──
  const commencementFormatted = formatDateString(
    quote.commencement_date,
    true,
    true,
    true
  )
  children.push(
    bodyText(
      `Dear ${quote.scheme_contact || scheme.contact_person || 'Sir/Madam'},`
    )
  )
  children.push(new Paragraph({ spacing: { after: 100 }, children: [] }))
  children.push(
    bodyText(
      `We are pleased to confirm that the group risk insurance scheme for ${quote.scheme_name} has been placed on risk ` +
        `with effect from ${commencementFormatted}. This letter serves as formal confirmation that cover is now active ` +
        `as per the accepted quotation.`
    )
  )
  children.push(new Paragraph({ spacing: { after: 100 }, children: [] }))

  // ── Cover Details table ──
  children.push(sectionHeading('Cover Details'))

  const labelCol = Math.round(CONTENT_WIDTH * 0.4)
  const valueCol = CONTENT_WIDTH - labelCol
  const colWidths = [labelCol, valueCol]

  const coverEndFormatted = formatDateString(
    quote.cover_end_date,
    true,
    true,
    true
  )

  const detailRows: [string, string][] = [
    ['Scheme Name', quote.scheme_name || ''],
    ['Quote Reference', `${quote.id}`],
    ['Risk Commencement Date', commencementFormatted],
    ['Cover End Date', coverEndFormatted],
    ['Industry', quote.industry || ''],
    ['Obligation Type', quote.obligation_type || ''],
    ['Number of Members', `${quote.member_data_count || 0}`],
    ['Currency', quote.currency || 'ZAR']
  ]

  children.push(
    new Table({
      width: { size: CONTENT_WIDTH, type: WidthType.DXA },
      columnWidths: colWidths,
      rows: detailRows.map(([k, v]) => kvRow(k, v, colWidths))
    })
  )

  children.push(new Paragraph({ spacing: { after: 200 }, children: [] }))

  // ── Benefits Summary table ──
  if (benefitSummary && benefitSummary.length > 0) {
    children.push(sectionHeading('Benefits Summary'))

    const benefitCol = Math.round(CONTENT_WIDTH * 0.6)
    const premiumCol = CONTENT_WIDTH - benefitCol
    const bColWidths = [benefitCol, premiumCol]

    // Header row
    const headerCells = ['Benefit', 'Annual Premium'].map(
      (label, i) =>
        new TableCell({
          width: { size: bColWidths[i], type: WidthType.DXA },
          shading: { fill: COLORS.primary, type: ShadingType.CLEAR },
          borders: thinBorders,
          margins: cellMargins,
          children: [
            new Paragraph({
              alignment: i === 0 ? AlignmentType.LEFT : AlignmentType.RIGHT,
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

    const dataRows = benefitSummary.map(
      (line: any, idx: number) =>
        new TableRow({
          children: [
            new TableCell({
              width: { size: bColWidths[0], type: WidthType.DXA },
              borders: thinBorders,
              margins: cellMargins,
              shading:
                idx % 2 === 1
                  ? { fill: COLORS.lightFill, type: ShadingType.CLEAR }
                  : undefined,
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: line.benefit,
                      font: FONT,
                      size: SIZE.body,
                      color: COLORS.dark
                    })
                  ]
                })
              ]
            }),
            new TableCell({
              width: { size: bColWidths[1], type: WidthType.DXA },
              borders: thinBorders,
              margins: cellMargins,
              shading:
                idx % 2 === 1
                  ? { fill: COLORS.lightFill, type: ShadingType.CLEAR }
                  : undefined,
              children: [
                new Paragraph({
                  alignment: AlignmentType.RIGHT,
                  children: [
                    new TextRun({
                      text: formatCurrency(
                        line.annual_premium,
                        quote.currency || 'ZAR'
                      ),
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

    // Total row
    const totalPremium = benefitSummary.reduce(
      (sum: number, l: any) => sum + (l.annual_premium || 0),
      0
    )
    const totalRow = new TableRow({
      children: [
        new TableCell({
          width: { size: bColWidths[0], type: WidthType.DXA },
          borders: thinBorders,
          margins: cellMargins,
          shading: { fill: COLORS.primary, type: ShadingType.CLEAR },
          children: [
            new Paragraph({
              children: [
                new TextRun({
                  text: 'Total Annual Premium',
                  font: FONT,
                  size: SIZE.body,
                  bold: true,
                  color: COLORS.white
                })
              ]
            })
          ]
        }),
        new TableCell({
          width: { size: bColWidths[1], type: WidthType.DXA },
          borders: thinBorders,
          margins: cellMargins,
          shading: { fill: COLORS.primary, type: ShadingType.CLEAR },
          children: [
            new Paragraph({
              alignment: AlignmentType.RIGHT,
              children: [
                new TextRun({
                  text: formatCurrency(totalPremium, quote.currency || 'ZAR'),
                  font: FONT,
                  size: SIZE.body,
                  bold: true,
                  color: COLORS.white
                })
              ]
            })
          ]
        })
      ]
    })

    children.push(
      new Table({
        width: { size: CONTENT_WIDTH, type: WidthType.DXA },
        columnWidths: bColWidths,
        rows: [
          new TableRow({ tableHeader: true, children: headerCells }),
          ...dataRows,
          totalRow
        ]
      })
    )

    children.push(new Paragraph({ spacing: { after: 200 }, children: [] }))
  }

  // ── Key Terms ──
  children.push(sectionHeading('Key Terms'))

  const termsColWidths = [labelCol, valueCol]
  const termsRows: [string, string][] = [
    [
      'Free Cover Limit',
      formatCurrency(quote.free_cover_limit || 0, quote.currency || 'ZAR')
    ],
    ['Normal Retirement Age', `${quote.normal_retirement_age || 'N/A'}`]
  ]

  children.push(
    new Table({
      width: { size: CONTENT_WIDTH, type: WidthType.DXA },
      columnWidths: termsColWidths,
      rows: termsRows.map(([k, v]) => kvRow(k, v, termsColWidths))
    })
  )

  children.push(new Paragraph({ spacing: { after: 200 }, children: [] }))

  // ── Closing paragraph ──
  const closingText =
    insurer.on_risk_letter_text ||
    `Should you have any queries regarding this confirmation of cover, please do not hesitate to contact us ` +
      `at ${insurer.email || 'our offices'}. We look forward to a continued relationship with ${quote.scheme_name}.`

  children.push(bodyText(closingText))
  children.push(new Paragraph({ spacing: { after: 300 }, children: [] }))
  children.push(bodyText('Yours faithfully,'))
  children.push(new Paragraph({ spacing: { after: 100 }, children: [] }))
  children.push(bodyText(insurer.name || '', { bold: true }))
  if (insurer.contact_person) {
    children.push(bodyText(insurer.contact_person))
  }

  return {
    properties: {
      page: {
        size: { width: A4.width, height: A4.height },
        margin: MARGINS
      }
    },
    headers: { default: makeHeader(quote.scheme_name || '') },
    footers: { default: makeFooter() },
    children
  }
}

// ---------------------------------------------------------------------------
// PDF generation
// ---------------------------------------------------------------------------

function generateOnRiskLetterPdfBlob(data: any): Blob {
  const {
    quote,
    scheme,
    insurer,
    letter,
    benefit_summary: benefitSummary
  } = data // eslint-disable-line camelcase

  // eslint-disable-next-line new-cap
  const doc = new jsPDF({ orientation: 'portrait', unit: 'mm', format: 'a4' })
  const pageWidth = doc.internal.pageSize.getWidth()
  const margin = 20
  let y = margin

  // Helper
  const addText = (
    text: string,
    x: number,
    yPos: number,
    options?: { fontSize?: number; fontStyle?: string; color?: number[] }
  ) => {
    doc.setFontSize(options?.fontSize || 10)
    doc.setFont('helvetica', options?.fontStyle || 'normal')
    if (options?.color) {
      doc.setTextColor(options.color[0], options.color[1], options.color[2])
    } else {
      doc.setTextColor(44, 62, 80)
    }
    doc.text(text, x, yPos)
  }

  // ── Insurer header ──
  addText(insurer.name || '', margin, y, {
    fontSize: 14,
    fontStyle: 'bold'
  })
  y += 6
  const addressParts = [
    insurer.address_line_1,
    insurer.address_line_2,
    insurer.address_line_3,
    [insurer.city, insurer.province, insurer.post_code]
      .filter(Boolean)
      .join(', '),
    insurer.telephone ? `Tel: ${insurer.telephone}` : '',
    insurer.email ? `Email: ${insurer.email}` : ''
  ].filter(Boolean)

  for (const line of addressParts) {
    addText(line, margin, y, { fontSize: 9 })
    y += 4
  }
  y += 6

  // ── Date ──
  addText(
    formatDateString(letter.letter_date || new Date(), true, true, true),
    margin,
    y,
    { fontSize: 10 }
  )
  y += 8

  // ── Addressee ──
  addText(quote.scheme_name || '', margin, y, {
    fontSize: 10,
    fontStyle: 'bold'
  })
  y += 5
  if (quote.scheme_contact || scheme.contact_person) {
    addText(quote.scheme_contact || scheme.contact_person, margin, y, {
      fontSize: 10
    })
    y += 5
  }
  y += 6

  // ── Subject ──
  addText('Confirmation of Cover — On Risk Letter', margin, y, {
    fontSize: 13,
    fontStyle: 'bold',
    color: [30, 58, 95]
  })
  y += 5
  if (letter.letter_reference) {
    addText(`Reference: ${letter.letter_reference}`, margin, y, {
      fontSize: 9,
      color: [88, 96, 105]
    })
    y += 6
  }
  y += 4

  // ── Opening paragraph ──
  const commFormatted = formatDateString(
    quote.commencement_date,
    true,
    true,
    true
  )
  const openingText =
    `We are pleased to confirm that the group risk insurance scheme for ${quote.scheme_name} has been placed on risk ` +
    `with effect from ${commFormatted}. This letter serves as formal confirmation that cover is now active ` +
    `as per the accepted quotation.`

  doc.setFontSize(10)
  doc.setFont('helvetica', 'normal')
  doc.setTextColor(44, 62, 80)
  const splitOpening = doc.splitTextToSize(openingText, pageWidth - 2 * margin)
  doc.text(splitOpening, margin, y)
  y += splitOpening.length * 5 + 6

  // ── Cover Details table ──
  addText('Cover Details', margin, y, {
    fontSize: 12,
    fontStyle: 'bold',
    color: [52, 73, 94]
  })
  y += 2

  const coverEndFormatted = formatDateString(
    quote.cover_end_date,
    true,
    true,
    true
  )
  const detailsBody = [
    ['Scheme Name', quote.scheme_name || ''],
    ['Quote Reference', `${quote.id}`],
    ['Risk Commencement Date', commFormatted],
    ['Cover End Date', coverEndFormatted],
    ['Industry', quote.industry || ''],
    ['Obligation Type', quote.obligation_type || ''],
    ['Number of Members', `${quote.member_data_count || 0}`],
    ['Currency', quote.currency || 'ZAR']
  ]

  ;(doc as any).autoTable({
    startY: y,
    head: [],
    body: detailsBody,
    theme: 'grid',
    margin: { left: margin, right: margin },
    styles: { fontSize: 9, cellPadding: 2 },
    columnStyles: {
      0: { fontStyle: 'bold', cellWidth: 55 }
    }
  })
  y = (doc as any).lastAutoTable.finalY + 8

  // ── Benefits Summary table ──
  if (benefitSummary && benefitSummary.length > 0) {
    addText('Benefits Summary', margin, y, {
      fontSize: 12,
      fontStyle: 'bold',
      color: [52, 73, 94]
    })
    y += 2

    const currency = quote.currency || 'ZAR'
    const symbol = currency === 'ZAR' ? 'R' : currency

    const benefitsBody = benefitSummary.map((line: any) => [
      line.benefit,
      `${symbol} ${line.annual_premium.toLocaleString('en-ZA', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
    ])

    const totalPremium = benefitSummary.reduce(
      (sum: number, l: any) => sum + (l.annual_premium || 0),
      0
    )
    benefitsBody.push([
      'Total Annual Premium',
      `${symbol} ${totalPremium.toLocaleString('en-ZA', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
    ])
    ;(doc as any).autoTable({
      startY: y,
      head: [['Benefit', 'Annual Premium']],
      body: benefitsBody,
      theme: 'grid',
      margin: { left: margin, right: margin },
      styles: { fontSize: 9, cellPadding: 2 },
      headStyles: {
        fillColor: [52, 73, 94],
        textColor: [255, 255, 255],
        fontStyle: 'bold'
      },
      columnStyles: {
        1: { halign: 'right' }
      },
      didParseCell: (hookData: any) => {
        // Bold the total row
        if (
          hookData.section === 'body' &&
          hookData.row.index === benefitsBody.length - 1
        ) {
          hookData.cell.styles.fontStyle = 'bold'
          hookData.cell.styles.fillColor = [52, 73, 94]
          hookData.cell.styles.textColor = [255, 255, 255]
        }
      }
    })
    y = (doc as any).lastAutoTable.finalY + 8
  }

  // ── Closing ──
  const closingText =
    insurer.on_risk_letter_text ||
    `Should you have any queries regarding this confirmation of cover, please do not hesitate to contact us.`

  doc.setFontSize(10)
  doc.setFont('helvetica', 'normal')
  doc.setTextColor(44, 62, 80)
  const splitClosing = doc.splitTextToSize(closingText, pageWidth - 2 * margin)
  doc.text(splitClosing, margin, y)
  y += splitClosing.length * 5 + 8

  addText('Yours faithfully,', margin, y)
  y += 8
  addText(insurer.name || '', margin, y, { fontStyle: 'bold' })
  y += 5
  if (insurer.contact_person) {
    addText(insurer.contact_person, margin, y)
  }

  // ── Footer ──
  const pageCount = (doc as any).internal.getNumberOfPages()
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i)
    doc.setFontSize(8)
    doc.setTextColor(150, 150, 150)
    doc.text(
      `Generated on ${formatDateString(new Date(), true, true, true)}`,
      margin,
      doc.internal.pageSize.getHeight() - 10
    )
    doc.text(
      `Page ${i} of ${pageCount}`,
      pageWidth - margin,
      doc.internal.pageSize.getHeight() - 10,
      { align: 'right' }
    )
  }

  return doc.output('blob') as unknown as Blob
}

// ---------------------------------------------------------------------------
// Composable
// ---------------------------------------------------------------------------

export function useOnRiskLetterGeneration() {
  const isGenerating = ref(false)
  const errorMessage = ref('')

  /**
   * Generate and download the On Risk letter as a Word (.docx) document.
   */
  async function generateOnRiskLetterDocx(data: any): Promise<void> {
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
        sections: [buildLetterSection(data)]
      })

      const blob = await Packer.toBlob(doc)
      const filename = `${data.quote.scheme_name}_On_Risk_Letter_${formatDateString(new Date(), true, true, true)}.docx`
      saveAs(blob, filename)
    } catch (err: any) {
      console.error('DOCX On Risk letter generation error:', err)
      errorMessage.value = `Error generating Word document: ${err.message}`
    } finally {
      isGenerating.value = false
    }
  }

  /**
   * Generate and download the On Risk letter as a PDF.
   */
  async function generateOnRiskLetterPdf(data: any): Promise<void> {
    isGenerating.value = true
    errorMessage.value = ''
    try {
      const blob = generateOnRiskLetterPdfBlob(data)
      const filename = `${data.quote.scheme_name}_On_Risk_Letter_${formatDateString(new Date(), true, true, true)}.pdf`
      saveAs(blob, filename)
    } catch (err: any) {
      console.error('PDF On Risk letter generation error:', err)
      errorMessage.value = `Error generating PDF: ${err.message}`
    } finally {
      isGenerating.value = false
    }
  }

  return {
    isGenerating,
    errorMessage,
    generateOnRiskLetterDocx,
    generateOnRiskLetterPdf
  }
}
