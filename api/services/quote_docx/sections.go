package quote_docx

import (
	"fmt"
	"strings"

	"api/models"
	"api/services"
)

// BuildCoverAndSummarySection builds the Cover and Summary section (Section 1)
// Portrait orientation. Header table (insurer info left, logo right), title, intro text, quote summary.
func BuildCoverAndSummarySection(quote models.GroupPricingQuote, summaries []models.MemberRatingResultSummary, insurer models.GroupPricingInsurerDetail) string {
	totals := CalculateQuoteTotals(summaries)
	summaryRows := BuildInitialInfoRows(quote, totals)

	var buf strings.Builder

	// Header table: insurer info (left) and logo (right)
	// This is a 2-column borderless table
	leftWidth := (ContentWidth * 6) / 10
	rightWidth := ContentWidth - leftWidth

	// Left cell: insurer details
	leftParas := []string{
		paragraphXML(ParagraphOptions{
			SpaceBefore: 60,
			SpaceAfter:  20,
		}, []string{
			runXML(RunOptions{
				Text:  insurer.Name,
				Bold:  true,
				Color: ColorDark,
				Size:  SizeSubheading,
			}),
		}),
		paragraphXML(ParagraphOptions{
			SpaceAfter: 10,
		}, []string{
			runXML(RunOptions{
				Text:  fmt.Sprintf("%s, %s", insurer.AddressLine1, insurer.AddressLine2),
				Color: ColorSecondary,
				Size:  SizeCaption,
			}),
		}),
	}

	if insurer.AddressLine3 != "" {
		leftParas = append(leftParas, paragraphXML(ParagraphOptions{
			SpaceAfter: 10,
		}, []string{
			runXML(RunOptions{
				Text:  insurer.AddressLine3,
				Color: ColorSecondary,
				Size:  SizeCaption,
			}),
		}))
	}

	leftParas = append(leftParas,
		paragraphXML(ParagraphOptions{
			SpaceAfter: 10,
		}, []string{
			runXML(RunOptions{
				Text:  fmt.Sprintf("%s, %s, %s", insurer.City, insurer.Province, insurer.PostCode),
				Color: ColorSecondary,
				Size:  SizeCaption,
			}),
		}),
		paragraphXML(ParagraphOptions{
			SpaceAfter: 10,
		}, []string{
			runXML(RunOptions{
				Text:  fmt.Sprintf("Tel: %s", insurer.Telephone),
				Color: ColorSecondary,
				Size:  SizeCaption,
			}),
		}),
		paragraphXML(ParagraphOptions{
			SpaceAfter: 10,
		}, []string{
			runXML(RunOptions{
				Text:  fmt.Sprintf("Email: %s", insurer.Email),
				Color: ColorSecondary,
				Size:  SizeCaption,
			}),
		}),
	)

	leftCell := cellXML(CellOptions{
		Width:       leftWidth,
		Shading:     ColorLightFill,
		HasBorders:  false,
	}, leftParas)

	// Right cell: logo or placeholder
	var rightParas []string
	if len(insurer.Logo) > 0 {
		// Use image run (simplified without full relationship setup here)
		rightParas = []string{
			paragraphXML(ParagraphOptions{
				Alignment:   "RIGHT",
				SpaceBefore: 60,
			}, []string{
				imageRunXML("rIdLogo", 210, 80, "Insurer logo"),
			}),
		}
	} else {
		rightParas = []string{
			paragraphXML(ParagraphOptions{
				Alignment: "RIGHT",
			}, []string{
				runXML(RunOptions{
					Text:   "[Company Logo]",
					Italic: true,
					Color:  "B4B4B4",
					Size:   SizeCaption,
				}),
			}),
		}
	}

	rightCell := cellXML(CellOptions{
		Width:        rightWidth,
		Shading:      ColorLightFill,
		HasBorders:   false,
		VerticalAlign: "CENTER",
	}, rightParas)

	// Header table
	headerRow := rowXML(RowOptions{}, []string{leftCell, rightCell})
	headerTable := tableXML(TableOptions{
		Width:        ContentWidth,
		ColumnWidths: []int{leftWidth, rightWidth},
	}, []string{headerRow})

	buf.WriteString(headerTable)
	buf.WriteString(spacerXML(200))

	// Title
	buf.WriteString(paragraphXML(ParagraphOptions{
		Alignment:   "CENTER",
		SpaceBefore: 100,
		SpaceAfter:  160,
	}, []string{
		runXML(RunOptions{
			Text:  "Group Risk Quotation",
			Bold:  true,
			Color: ColorDark,
			Size:  SizeTitle,
		}),
	}))

	// Intro text
	introText := insurer.IntroductoryText
	if introText == "" {
		introText = "We are pleased to submit for your consideration the quotation you requested for the above scheme."
	}
	buf.WriteString(bodyTextXML(introText, false))
	buf.WriteString(spacerXML(80))

	// Quote summary table
	buf.WriteString(keyValueTableXML(summaryRows, ContentWidth))

	// Section properties (portrait)
	buf.WriteString(sectionPropsXML("PORTRAIT", MarginTop, MarginBottom, MarginLeft, MarginRight))

	return buf.String()
}

// BuildPremiumSummarySection builds the Premium Summary section (Section 2)
// Landscape orientation. Conditionally includes Premium Summary table if non-funeral benefits exist.
func BuildPremiumSummarySection(quote models.GroupPricingQuote, summaries []models.MemberRatingResultSummary) string {
	hasNonFuneral := HasAnyNonFuneralBenefits(summaries)
	var buf strings.Builder

	if hasNonFuneral {
		buf.WriteString(bodyTextXML("The following table provides a summary of the benefits and premiums for the scheme:", false))
		buf.WriteString(sectionHeadingXML("Premium Summary"))

		premiumRows := BuildPremiumSummaryRows(summaries)
		labels := []string{"Category", "No of Lives", "Total Salary", "Total Sum Assured", "Annual Premium", "% Salary"}

		cw := LandscapeContentWidth
		colWidths := []int{
			(cw * 2) / 10,
			(cw * 12) / 100,
			(cw * 2) / 10,
			(cw * 2) / 10,
			(cw * 18) / 100,
			(cw * 1) / 10,
		}
		colWidths[5] = cw - (colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3] + colWidths[4])

		tableRows := []string{headerRowXML(labels, colWidths)}
		for i, row := range premiumRows {
			isTotal := row.Category == "Total"
			fillColor := ""
			if isTotal {
				fillColor = ColorLightFill
			} else if i%2 == 1 {
				fillColor = "FCFDFE"
			}

			alignments := []string{"LEFT", "CENTER", "RIGHT", "RIGHT", "RIGHT", "CENTER"}
			values := []string{row.Category, row.MemberCount, row.TotalSalary, row.TotalSumAssured, row.AnnualPremium, row.PercentSalary}
			tableRows = append(tableRows, dataRowXML(values, colWidths, alignments, isTotal, fillColor))
		}

		buf.WriteString(tableXML(TableOptions{
			Width:        cw,
			ColumnWidths: colWidths,
		}, tableRows))
	}

	// Group Funeral table (always shown)
	buf.WriteString(spacerXML(200))
	buf.WriteString(sectionHeadingXML("Group Funeral"))

	funeralRows := BuildGroupFuneralRows(summaries)
	fLabels := []string{"Category", "No of Lives", "Monthly Premium", "Annual Premium", "Total Annual Premium"}

	fcw := LandscapeContentWidth
	fColWidths := []int{
		(fcw * 22) / 100,
		(fcw * 14) / 100,
		(fcw * 2) / 10,
		(fcw * 2) / 10,
		(fcw * 24) / 100,
	}
	fColWidths[4] = fcw - (fColWidths[0] + fColWidths[1] + fColWidths[2] + fColWidths[3])

	fTableRows := []string{headerRowXML(fLabels, fColWidths)}
	for i, row := range funeralRows {
		isTotal := row.Category == "Total"
		fillColor := ""
		if isTotal {
			fillColor = ColorLightFill
		} else if i%2 == 1 {
			fillColor = "FCFDFE"
		}

		alignments := []string{"LEFT", "CENTER", "RIGHT", "RIGHT", "RIGHT"}
		values := []string{row.Category, row.MemberCount, row.MonthlyPremium, row.AnnualPremium, row.TotalAnnualPremium}
		fTableRows = append(fTableRows, dataRowXML(values, fColWidths, alignments, isTotal, fillColor))
	}

	buf.WriteString(tableXML(TableOptions{
		Width:        fcw,
		ColumnWidths: fColWidths,
	}, fTableRows))

	// Section properties (landscape)
	buf.WriteString(sectionPropsXML("LANDSCAPE", MarginTop, MarginBottom, MarginLeft, MarginRight))

	return buf.String()
}

// BuildPremiumBreakdownSection builds the Premium Breakdown section (Section 3)
// Landscape orientation. Per category: heading, benefit breakdown table, funeral sub-table.
func BuildPremiumBreakdownSection(quote models.GroupPricingQuote, summaries []models.MemberRatingResultSummary, titles BenefitTitles) string {
	hasNonFuneral := HasAnyNonFuneralBenefits(summaries)
	var buf strings.Builder

	buf.WriteString(sectionHeadingXML("Premium Breakdown"))

	cw := LandscapeContentWidth

	for _, item := range summaries {
		if hasNonFuneral {
			buf.WriteString(categoryHeadingXML(fmt.Sprintf("%s Category", item.Category)))
		}

		// Benefit breakdown table
		if CategoryHasNonFuneralBenefits(item) {
			breakdownRows := BuildPremiumBreakdownRows(item, titles)
			bLabels := []string{"Benefit", "Total Sum Assured", "Annual Premium", "% Salary"}
			bColWidths := []int{
				(cw * 3) / 10,
				(cw * 25) / 100,
				(cw * 25) / 100,
				(cw * 2) / 10,
			}
			bColWidths[3] = cw - (bColWidths[0] + bColWidths[1] + bColWidths[2])

			bTableRows := []string{headerRowXML(bLabels, bColWidths)}
			for i, row := range breakdownRows {
				fillColor := ""
				if i%2 == 1 {
					fillColor = "FDFDFE"
				}
				alignments := []string{"LEFT", "RIGHT", "RIGHT", "CENTER"}
				values := []string{row.Benefit, row.TotalSumAssured, row.AnnualPremium, row.PercentSalary}
				bTableRows = append(bTableRows, dataRowXML(values, bColWidths, alignments, false, fillColor))
			}

			buf.WriteString(tableXML(TableOptions{
				Width:        cw,
				ColumnWidths: bColWidths,
			}, bTableRows))
		}

		// Group Funeral sub-table for this category
		buf.WriteString(paragraphXML(ParagraphOptions{
			SpaceBefore: 120,
			SpaceAfter:  40,
		}, []string{
			runXML(RunOptions{
				Text:  fmt.Sprintf("%s - Group Funeral", item.Category),
				Bold:  true,
				Color: ColorDark,
				Size:  SizeBody,
			}),
		}))

		funeralKV := BuildGroupFuneralBreakdownRows(item)
		buf.WriteString(keyValueTableXML(funeralKV, (cw*5)/10))

		buf.WriteString(spacerXML(100))
	}

	// Section properties (landscape)
	buf.WriteString(sectionPropsXML("LANDSCAPE", MarginTop, MarginBottom, MarginLeft, MarginRight))

	return buf.String()
}

// BuildBenefitsDefinitionsSection builds the Benefits and Definitions section (Section 4)
// Landscape orientation. Per category: common benefits, 7-column definition table, funeral coverage, educator benefits.
func BuildBenefitsDefinitionsSection(quote models.GroupPricingQuote, summaries []models.MemberRatingResultSummary, educatorBenefits []interface{}, titles BenefitTitles) string {
	hasNonFuneral := HasAnyNonFuneralBenefits(summaries)
	var buf strings.Builder

	buf.WriteString(sectionHeadingXML("Benefits and Definitions of the Cover"))

	cw := LandscapeContentWidth

	for _, cat := range quote.SchemeCategories {
		buf.WriteString(categoryHeadingXML(cat.SchemeCategory))

		// Common benefits key-value table
		if hasNonFuneral {
			commonRows := BuildCategoryCommonBenefitRows(cat, quote)
			buf.WriteString(keyValueTableXML(commonRows, (cw*55)/100))
			buf.WriteString(spacerXML(60))
		}

		// 7-column benefit definitions table
		if hasNonFuneral {
			defRows := BuildBenefitDefinitionRows(cat, quote, titles)
			defLabels := []string{"Benefit", "Salary Multiple", "Benefit Structure", "Waiting Period", "Deferred Period", "Cover Definition", "Risk Type"}
			dColWidths := []int{
				(cw * 16) / 100,
				(cw * 12) / 100,
				(cw * 13) / 100,
				(cw * 12) / 100,
				(cw * 12) / 100,
				(cw * 18) / 100,
				(cw * 17) / 100,
			}
			dColWidths[6] = cw - (dColWidths[0] + dColWidths[1] + dColWidths[2] + dColWidths[3] + dColWidths[4] + dColWidths[5])

			dTableRows := []string{headerRowXML(defLabels, dColWidths)}
			for i, row := range defRows {
				fillColor := ""
				if i%2 == 1 {
					fillColor = "FCFDFE"
				}
				alignments := []string{"LEFT", "CENTER", "CENTER", "CENTER", "CENTER", "CENTER", "CENTER"}
				values := []string{row.Benefit, row.SalaryMultiple, row.BenefitStructure, row.WaitingPeriod, row.DeferredPeriod, row.CoverDefinition, row.RiskType}
				dTableRows = append(dTableRows, dataRowXML(values, dColWidths, alignments, false, fillColor))
			}

			buf.WriteString(tableXML(TableOptions{
				Width:        cw,
				ColumnWidths: dColWidths,
			}, dTableRows))
			buf.WriteString(spacerXML(80))
		}

		// Group Funeral coverage table
		buf.WriteString(paragraphXML(ParagraphOptions{
			SpaceBefore: 80,
			SpaceAfter:  40,
		}, []string{
			runXML(RunOptions{
				Text:  fmt.Sprintf("%s - Group Funeral", cat.SchemeCategory),
				Bold:  true,
				Color: ColorDark,
				Size:  SizeBody,
			}),
		}))

		funeralCovRows := BuildFuneralCoverageRows(cat)
		fLabels := []string{"Member", "Sum Assured", "Maximum Number Covered"}
		fcw := (cw * 55) / 100
		fColWidths := []int{
			(fcw * 4) / 10,
			(fcw * 3) / 10,
			(fcw * 3) / 10,
		}
		fColWidths[2] = fcw - fColWidths[0] - fColWidths[1]

		fRows := []string{headerRowXML(fLabels, fColWidths)}
		for i, row := range funeralCovRows {
			fillColor := ""
			if i%2 == 1 {
				fillColor = "FCFDFE"
			}
			alignments := []string{"LEFT", "RIGHT", "CENTER"}
			values := []string{
				row.Member,
				fmt.Sprintf("%.0f", row.SumAssured),
				fmt.Sprintf("%.0f", row.MaxCovered),
			}
			fRows = append(fRows, dataRowXML(values, fColWidths, alignments, false, fillColor))
		}

		buf.WriteString(tableXML(TableOptions{
			Width:        fcw,
			ColumnWidths: fColWidths,
		}, fRows))

		buf.WriteString(spacerXML(200))
	}

	// Section properties (landscape)
	buf.WriteString(sectionPropsXML("LANDSCAPE", MarginTop, MarginBottom, MarginLeft, MarginRight))

	return buf.String()
}

// BuildProvisionsSection builds the Provisions section (Section 5)
// Portrait orientation. Provisions text + closing block.
func BuildProvisionsSection(quote models.GroupPricingQuote, insurer models.GroupPricingInsurerDetail) string {
	var buf strings.Builder

	buf.WriteString(sectionHeadingXML("Underwriting and General Provisions"))

	provisionsText := insurer.GeneralProvisionsText
	if provisionsText == "" {
		provisionsText = "Please refer to the terms and conditions of the policy document for full details on underwriting and general provisions."
	}
	buf.WriteString(bodyTextXML(provisionsText, false))

	// Closing block
	buf.WriteString(spacerXML(400))

	// Thank you paragraphs with light grey background
	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill: ColorLightFill,
	}, []string{
		runXML(RunOptions{
			Text:   "Thank you for considering our quotation.",
			Italic: true,
			Color:  ColorSecondary,
			Size:   SizeBody,
		}),
	}))

	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill: ColorLightFill,
	}, []string{
		runXML(RunOptions{
			Text:   "We look forward to the opportunity to serve your insurance needs.",
			Italic: true,
			Color:  ColorSecondary,
			Size:   SizeBody,
		}),
	}))

	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:        ColorLightFill,
		SpaceBefore: 80,
	}, []string{
		runXML(RunOptions{
			Text:  insurer.Name,
			Bold:  true,
			Color: ColorSecondary,
			Size:  SizeBody,
		}),
	}))

	// Section properties (portrait)
	buf.WriteString(sectionPropsXML("PORTRAIT", MarginTop, MarginBottom, MarginLeft, MarginRight))

	return buf.String()
}

// BuildAcceptanceFormSection builds the Acceptance Form section (Section 6)
// Portrait orientation. Navy header, policy details form, signature boxes, POPIA consent, office use table.
func BuildAcceptanceFormSection(quote models.GroupPricingQuote) string {
	cw := ContentWidth
	var buf strings.Builder

	// Navy header bar
	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill: ColorNavy,
	}, []string{
		runXML(RunOptions{
			Text:  "ACCEPTANCE OF QUOTATION",
			Bold:  true,
			Color: ColorWhite,
			Size:  SizeTitle,
		}),
	}))

	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:       ColorNavy,
		SpaceAfter: 100,
	}, []string{
		runXML(RunOptions{
			Text:  "POPIA Compliant",
			Color: "E6E6E6",
			Size:  SizeSmall,
		}),
	}))

	// Policy Details section
	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:        "FAFBFC",
		SpaceBefore: 200,
		SpaceAfter:  80,
	}, []string{
		runXML(RunOptions{
			Text:  "POLICY DETAILS",
			Bold:  true,
			Color: ColorNavy,
			Size:  14,
		}),
	}))

	// Form table helper - creates a simple form row
	policyTable := buildFormTable(cw, [][]string{
		{"EMPLOYER / SCHEME NAME", "QUOTE NUMBER"},
		{"DATE OF QUOTE (DD/MM/YYYY)", "COMMENCEMENT DATE (DD/MM/YYYY)"},
	})
	buf.WriteString(policyTable)

	// Profile notice — tolerance percentage is sourced from the singleton
	// GroupPricingSetting row so the wording stays in sync with the value
	// configured in the Risk Watchlist Thresholds metadata panel.
	buf.WriteString(spacerXML(160))
	tolerancePct := services.GetRiskProfileVariationTolerancePct()
	buf.WriteString(bodyTextXML(
		fmt.Sprintf("If the member data profile at the quotation implementation date differ by %s%% or more from that on which the quotation was based, we reserve the right to revise the rates and Automatic Acceptance Limit. The Employer/Scheme will be notified accordingly and must provide acceptance before implementation proceeds.", formatTolerancePct(tolerancePct)),
		false,
	))
	buf.WriteString(bodyTextXML(
		"By signing this quotation, the Employer/Scheme acknowledges that they have read, understood, and agree to be bound by all the terms and conditions of this quotation.",
		false,
	))

	// Employer Authorisation
	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:        "FAFBFC",
		SpaceBefore: 200,
		SpaceAfter:  80,
	}, []string{
		runXML(RunOptions{
			Text:  "EMPLOYER AUTHORISATION",
			Bold:  true,
			Color: ColorNavy,
			Size:  14,
		}),
	}))

	// Signature box (simplified)
	buf.WriteString(buildSignatureBox("Duly Authorised Signatory"))
	buf.WriteString(spacerXML(60))

	// Signature details form
	sigTable := buildFormTable(cw, [][]string{
		{"FULL NAME", "CAPACITY (e.g., Director)", "DATE (DD/MM/YYYY)"},
	})
	buf.WriteString(sigTable)

	// Intermediary Details
	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:        "FAFBFC",
		SpaceBefore: 200,
		SpaceAfter:  80,
	}, []string{
		runXML(RunOptions{
			Text:  "INTERMEDIARY DETAILS",
			Bold:  true,
			Color: ColorNavy,
			Size:  14,
		}),
	}))

	buf.WriteString(buildSignatureBox("Intermediary / FAIS Representative"))
	buf.WriteString(spacerXML(60))

	intTable := buildFormTable(cw, [][]string{
		{"FULL NAME", "FAIS REG NO.", "DATE (DD/MM/YYYY)"},
	})
	buf.WriteString(intTable)

	// POPIA Consent
	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:        ColorLightBlue,
		SpaceBefore: 200,
		SpaceAfter:  40,
	}, []string{
		runXML(RunOptions{
			Text:  "POPIA Consent & Data Protection",
			Bold:  true,
			Color: "0366D6",
			Size:  14,
		}),
	}))

	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:       ColorLightBlue,
		SpaceAfter: 60,
	}, []string{
		runXML(RunOptions{
			Text:  "In terms of the Protection of Personal Information Act 4 of 2013 (POPIA), the Employer consents to the processing of personal information of employees and scheme members for the purpose of underwriting, administering, and processing claims under this Group Risk policy. Information will be processed lawfully, minimally, and only for the specific purpose stated.",
			Color: ColorDark,
			Size:  14,
		}),
	}))

	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:       ColorLightBlue,
		SpaceAfter: 80,
	}, []string{
		runXML(RunOptions{
			Text:  "☐ I confirm that the Employer has obtained necessary consent from data subjects (employees/members) for the processing of their personal information as required by POPIA, and warrants that all information provided is true and complete.",
			Color: ColorDark,
			Size:  13,
		}),
	}))

	// For Office Use Only
	buf.WriteString(paragraphXML(ParagraphOptions{
		Fill:        ColorLightOrange,
		SpaceBefore: 200,
		SpaceAfter:  80,
	}, []string{
		runXML(RunOptions{
			Text:  "FOR OFFICE USE ONLY",
			Bold:  true,
			Color: ColorOrange,
			Size:  14,
		}),
	}))

	// Office grid table
	colW := cw / 3
	lastColW := cw - colW*2
	officeBorder := "FFB47C"

	row1 := buildOfficeGridRow([]string{"RECEIVED BY", "DATE RECEIVED", "POLICY NUMBER"}, colW, lastColW, officeBorder)
	row2 := buildOfficeGridRow([]string{"UNDERWRITER", "APPROVED BY", "DATE APPROVED"}, colW, lastColW, officeBorder)

	buf.WriteString(tableXML(TableOptions{
		Width:        cw,
		ColumnWidths: []int{colW, colW, lastColW},
	}, []string{row1, row2}))

	// Section properties (portrait)
	buf.WriteString(sectionPropsXML("PORTRAIT", MarginTop, MarginBottom, MarginLeft, MarginRight))

	return buf.String()
}

// Helper: Build form table with variable columns
func buildFormTable(width int, rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	numCols := len(rows[0])
	colWidth := width / numCols

	var tableRows []string
	for _, row := range rows {
		cells := []string{}
		for i, label := range row {
			width := colWidth
			if i == numCols-1 {
				width = width - colWidth*(numCols-1)
			}

			para := paragraphXML(ParagraphOptions{
				SpaceAfter: 20,
			}, []string{
				runXML(RunOptions{
					Text:  label,
					Color: ColorMediumGray,
					Size:  12,
				}),
			})

			// Input line
			inputLine := paragraphXML(ParagraphOptions{
				SpaceBefore: 80,
				SpaceAfter:  40,
				BorderBottom: true,
				BorderColor: ColorInputLine,
			}, []string{
				runXML(RunOptions{
					Text:  " ",
					Size:  SizeBody,
				}),
			})

			cells = append(cells, cellXML(CellOptions{
				Width:       width,
				HasBorders:  false,
			}, []string{para, inputLine}))
		}

		tableRows = append(tableRows, rowXML(RowOptions{}, cells))
	}

	return tableXML(TableOptions{
		Width:        width,
		ColumnWidths: make([]int, numCols),
	}, tableRows)
}

// Helper: Build signature box
func buildSignatureBox(title string) string {
	boxWidth := ContentWidth - 120
	para1 := paragraphXML(ParagraphOptions{
		SpaceAfter: 20,
	}, []string{
		runXML(RunOptions{
			Text:  title,
			Bold:  true,
			Color: ColorNavy,
			Size:  14,
		}),
	})

	// Space for signature
	para2 := paragraphXML(ParagraphOptions{
		SpaceBefore: 300,
		SpaceAfter:  0,
	}, []string{})

	// Signature line
	para3 := paragraphXML(ParagraphOptions{
		BorderBottom: true,
		BorderColor:  ColorInputLine,
	}, []string{
		runXML(RunOptions{
			Text:  "Sign here",
			Color: "969696",
			Size:  12,
		}),
	})

	cell := cellXML(CellOptions{
		Width:       boxWidth,
		Shading:     ColorWhite,
		HasBorders:  true,
		BorderColor: "000000",
	}, []string{para1, para2, para3})

	return tableXML(TableOptions{
		Width:        boxWidth,
		ColumnWidths: []int{boxWidth},
	}, []string{rowXML(RowOptions{}, []string{cell})})
}

// Helper: Build office grid row
func buildOfficeGridRow(labels []string, colW, lastColW int, borderColor string) string {
	cells := []string{}
	for i, label := range labels {
		width := colW
		if i == len(labels)-1 {
			width = lastColW
		}

		para := paragraphXML(ParagraphOptions{}, []string{
			runXML(RunOptions{
				Text:  label,
				Color: ColorMediumGray,
				Size:  12,
			}),
		})

		cells = append(cells, cellXML(CellOptions{
			Width:       width,
			Shading:     ColorLightOrange,
			HasBorders:  true,
			BorderColor: borderColor,
		}, []string{para}))
	}

	return rowXML(RowOptions{}, cells)
}
