export const groupPricing = [
  {
    data_variable: 'member_distribution_free_cover_limit',
    data_type: 'number',
    data_description:
      'The free cover limit is the maximum coverage an individual can receive without additional medical information or underwriting. It is calculated based on the underlying distribution of members’ sum assured.',
    data_source:
      'Min( (FreeCoverLimitScalingFactor ) * SQRT(NumberOfMembers) * AverageUncappedSumAssured,\n\t\t\t\t\t Percentile(MemberData,FreeCoverLimitPercentile) )',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'free_cover_limit_scaling_factor',
    data_type: 'number',
    data_description:
      'Allows the free cover limit to be scaled appropriately based on the group size and the average sum assured, helping to balance risk while providing a coverage limit. The final free cover limit is constrained by the Percentile(MemberData, FreeCoverLimitPercentile), ensuring it stays within acceptable limits. Ref. MemeberDistributionFreeCoverLimit formula for details',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'free_cover_limit_percentile',
    data_type: 'number',
    data_description:
      'Refers to the percentile [0,1] applied to the MemberData (e.g., the sum assured or coverage amounts of the members) to determine the limit.It ensures that the free cover limit is not set too high, regardless of the scaling factor and other parameters. It accounts for the spread or distribution of member data, ensuring that the free cover limit reflects the overall group characteristics while avoiding overly generous limits that could pose higher risk.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'free_cover_limit_nearest_multiple',
    data_type: 'number',
    data_description:
      'A specified multiple to which the MemberDistributionFreeCoverLimit is rounded. This multiple may be a fixed amount—such as 10,000, 50,000, or 100,000—depending on the policy or applicable calculation rules.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'fcl_maximum_cover_scaling_factor',
    data_type: 'number',
    data_description:
      'Multiplier applied to the largest member sum assured to derive an upper bound on the free cover limit when the Statistical Outlier method is in use. The final free cover limit is then capped at this product, alongside the scaled mean salary and the log-normal upper bound on sum assured. Only consulted when the system-wide free cover limit method is set to Statistical Outlier; ignored under the Percentile method.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'maximum_allowed_fcl',
    data_type: 'number',
    data_description:
      'Underwriting ceiling for user-set free cover limits, configured per risk rate. When a quote-level free cover limit is enforced, the system uses it as-is unless it exceeds this ceiling by more than the configured override tolerance, in which case the ceiling is used. A value of 0 means no ceiling is configured and quote-level overrides pass through unchanged. Applies under both the Percentile and Statistical Outlier calculation methods.',
    data_source: 'Restrictions Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'fcl_override_tolerance',
    data_type: 'number',
    data_description:
      'Fractional headroom (0–1) allowed above the maximum allowed free cover limit before a quote-level override is clamped. Configured system-wide alongside the free cover limit calculation method on the metadata configuration screen. A value of 0.2 means quote-level overrides up to 20% above the ceiling are honoured as-is; anything beyond is clamped to the ceiling.',
    data_source: 'Group Pricing Settings (singleton metadata)',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'free_cover_limit',
    data_type: 'number',
    data_description:
      'The free cover limit applied to the scheme category. If no quote-level free cover limit is enforced, equals the member-distribution-derived value. If enforced, equals the quote-level value, clamped to the maximum allowed free cover limit when the quote-level value exceeds it by more than the configured override tolerance.',
    data_source:
      'If quote-level FreeCoverLimit = 0: MemberDistributionFreeCoverLimit.\nIf quote-level FreeCoverLimit > 0 and MaximumAllowedFCL > 0 and FreeCoverLimit > (1 + OverrideTolerance) * MaximumAllowedFCL: MaximumAllowedFCL.\nOtherwise: quote-level FreeCoverLimit.',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'commencement_date',
    data_type: 'number',
    data_description: 'The date on which the group risk cover starts',
    data_source: ' Group Quote',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'entry_date',
    data_type: 'number',
    data_description: 'The date on which a member entered the group scheme.',
    data_source: ' Original Member Data or Scheme Maintenance',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'exit_date',
    data_type: 'number',
    data_description: 'The date on which a member exits the group scheme.',
    data_source: ' Scheme Maintenance',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'education_level',
    data_type: 'number',
    data_description:
      'Education Levels 1 - Pre-School, 2- Primary 3- Secondary 4 - Tertiary respectively',
    data_source: ' Group Pricing Educator Structure',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_coverage_period',
    data_type: 'number',
    data_description:
      'The number of years over which the education benefit is provided under each respective education level. This duration may vary by level (Pre-School, Primary, Secondary, Tertiary) and is defined according to the scheme’s terms and conditions. Ref. Group Pricing Educator Structure Table',
    data_source: ' Group Pricing Educator Structure',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_tuition_per_year',
    data_type: 'number',
    data_description:
      'The maximum educator benefit amount available per year under each respective education level. For detailed information, please refer to the Group Pricing Educator Structure Table.',
    data_source: ' Group Pricing Educator Structure',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_book_allowance_proportion',
    data_type: 'number',
    data_description:
      'The maximum educator book allowance is calculated as a proportion of the respective education levels maximum tuition per year. For detailed information, please refer to the Group Pricing Educator Structure Table.',
    data_source: ' Group Pricing Educator Structure',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_book_allowance_amount',
    data_type: 'number',
    data_description:
      'The maximum educator book allowance amount per year. For detailed information, please refer to the Group Pricing Educator Structure Table.',
    data_source: ' Group Pricing Educator Structure',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_book_allowance',
    data_type: 'number',
    data_description:
      'The maximum educator book allowance per year is determined by the higher of the amount calculated from the MaxBookAllowanceProportion and the MaxBookAllowanceAmount. For detailed information, please refer to the Group Pricing Educator Structure Table.',
    data_source:
      'Max( MaxBookAllowanceAmount, MaxBookAllowanceProportion * MaxTuitionPerYear)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'max_accommodation_allowance_proportion',
    data_type: 'number',
    data_description:
      'The maximum educator accommodation allowance is calculated as a proportion of the respective education levels maximum tuition per year. For detailed information, please refer to the Group Pricing Educator Structure Table.',
    data_source: ' Group Pricing Educator Structure',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_accommodation_allowance_amount',
    data_type: 'number',
    data_description:
      'The maximum educator accommodation allowance amount per year. For detailed information, please refer to the Group Pricing Educator Structure Table.',
    data_source: ' Group Pricing Educator Structure',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_accommodation_allowance',
    data_type: 'number',
    data_description:
      'The maximum educator accommodation allowance per year is determined by the higher of the amount calculated from the MaxAccommodationAllowanceProportion and the MaxAccommodationAllowanceAmount. For detailed information, please refer to the Group Pricing Educator Structure Table.',
    data_source:
      'Max( MaxAccomodationAllowanceAmount, MaxAccommodationAllowanceProportion * MaxTuitionPerYear)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'grade0_sum_assured',
    data_type: 'number',
    data_description: 'Represents the total sum assured exposure for Grade 0',
    data_source:
      'Grade0MaxTuitionPerYear * Grade0MaxCoverageYears *\n\t(1 + MaxBookAllowanceProportion+ MaxAccommodationAllowanceProportion)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'grade1_7_sum_assured',
    data_type: 'number',
    data_description:
      'Represents the total sum assured exposure for Grade 1 to Grade 7',
    data_source:
      'Grade17MaxTuitionPerYear * Grade17MaxCoverageYears *\n\t(1 + MaxBookAllowanceProportion+ MaxAccommodationAllowanceProportion)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'grade8_12_sum_assured',
    data_type: 'number',
    data_description:
      'Represents the total sum assured exposure for Grade 8 to Grade 12',
    data_source:
      'Grade812MaxTuitionPerYear * Grade812MaxCoverageYears *\n\t(1 + MaxBookAllowanceProportion+ MaxAccommodationAllowanceProportion)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'tertiary_sum_assured',
    data_type: 'number',
    data_description: 'Represents the total sum assured exposure for tertiary',
    data_source:
      'TertiaryMaxTuitionPerYear * TertiaryMaxCoverageYears *\n\t(1 + MaxBookAllowanceProportion+ MaxAccommodationAllowanceProportion)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'gla_educator_risk_premium',
    data_type: 'number',
    data_description:
      'GLA component of the educator risk premium. Present when the GLA educator benefit is enabled on the scheme category.',
    data_source:
      '( Grade0SumAssured*Grade0RiskRate + Grade17SumAssured*Grade17RiskRate\n\t\t + Grade812SumAssured*Grade812RiskRate + TertiarySumAssured*TertiaryRiskRate ) *\n\t\tLoadedGlaRate',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'ptd_educator_risk_premium',
    data_type: 'number',
    data_description:
      'PTD component of the educator risk premium. Present when the PTD educator benefit is enabled on the scheme category.',
    data_source:
      '( Grade0SumAssured*Grade0RiskRate + Grade17SumAssured*Grade17RiskRate\n\t\t + Grade812SumAssured*Grade812RiskRate + TertiarySumAssured*TertiaryRiskRate ) *\n\t\tLoadedPtdRate',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'spouse_age_gap',
    data_type: 'number',
    data_description:
      'Age gap between the main member and their spouse. Used to estimate the spouse`s age from the main member`s age, subject to min_age and max_age constraints',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'min_age',
    data_type: 'number',
    data_description:
      'Minimum age constraint used as a lower bound when estimating the spouse`s age.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'max_age',
    data_type: 'number',
    data_description:
      'Maximum age constraint used as an upper bound when estimating the spouse`s age.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'is_lumpsum_reins_gla_dependent',
    data_type: 'number',
    data_description:
      'Boolean variable indicating whether the ceded proportion of other lump sum components is based on the underlying GLA reinsurance cession ratios. A value of 1 indicates True; 0 indicates False.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'gla_terminal_illness_loading_rate',
    data_type: 'number',
    data_description:
      'A rate between 0 and 1 representing a proportional loading over the base GLA rate to reflect the additional terminal illness risk, if included under the GLA benefit configuration.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'premium_rates_guaranteed_period_months',
    data_type: 'number',
    data_description:
      'The number of months during which premium rates are guaranteed and cannot be subject to review.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'quote_validity_period_months',
    data_type: 'number',
    data_description:
      'Duration in months for which the quote remains valid, starting from the quote date',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'annual_expense_amount',
    data_type: 'number',
    data_description:
      'Annual expense amount per member, in addition to the expense loading.',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'full_credibility_threshold',
    data_type: 'number',
    data_description:
      'The full credibility threshold represents the minimum amount of exposure required for a group`s experience to be considered fully credible (i.e., 100% weight in calculations)',
    data_source: 'Group Pricing Parameter Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'credibility',
    data_type: 'number',
    data_description:
      'Refers to the degree of confidence or weight given to the experience data',
    data_source: 'Min( Sqrt( WeightedLifeYears/FullCredibilityThreshold ), 1 )',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'blended_gla_rate',
    data_type: 'number',
    data_description:
      'The blended Group Life Assurance (GLA) rate is calculated to provide a balanced view between the actual claims experience of a specific scheme and the insurer’s base rates or community rates, with each component weighted according to its credibility',
    data_source:
      'GlaTheoreticalRate *(1- credibility) + credibility * AnnualExperienceWeightedRate',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'annual_experience_weighted_rate',
    data_type: 'number',
    data_description:
      'is the observed cost of claims for a group during a specific experience period, expressed as a rate relative to the insured exposure (for example insured salary, sum assured, or member count depending on the product). It reflects the group’s actual risk experience before credibility adjustments.For experience rating across multiple years, the calculated experience rates for each year are weighted using the defined annual weightings to produce a weighted experience rate that reflects the relative importance of each experience period.',
    data_source:
      '+= ((claimsDataPoint.GlaClaimsAmount / (claimsDataPoint.TotalGlaSumAssured))/timePeriodYears ) * claimsDataPoint.Weighting',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'gla_experience_adjustment',
    data_type: 'number',
    data_description:
      'The GLA Experience Adjustment Factor is a scaling factor applied to the community rate or base rate to reflect the actual claims experience of a specific group scheme',
    data_source: 'blended_gla_rate/GlaTheoreticalRate',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'gla_theoretical_rate',
    data_type: 'number',
    data_description:
      'is the expected or manual premium rate for Group Life Assurance (GLA) calculated from the insurer`s pricing basis rather than the group`s claims experience. It represents the rate that would apply to the group based purely on its risk characteristics (such as demographics, benefit structure, and underwriting assumptions) before experience rating adjustments are applied.Typically derived from pricing tables and risk factors',
    data_source:
      'expected risk rate for the exposure data in question, read from the basis table',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'mannually_added_credibility',
    data_type: 'number',
    data_description:
      'A manually added credibility factor. It is based on actuarial judgement, taking into account the calculated credibility as well as credibility levels historically applied to similar schemes with comparable exposure and experience characteristics',
    data_source: 'manual input within member rating results',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'base_rate',
    data_type: 'number',
    data_description:
      'The starting rate per benefit, derived by grossing up the underlying mortality / morbidity rate (Qx, read from the basis table by age, gender, occupational class) for industry and region risk. Applies per benefit (GLA, PTD, CI, TTD, PHI, Funeral, etc.).',
    data_source: 'Qx * (1 + IndustryLoading + RegionLoading)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'loaded_rate',
    data_type: 'number',
    data_description:
      'The Base Rate after the per-member loadings on top of base risk are added. Loadings are summed (not compounded). SpecialLoadings sums the benefit-specific extras such as Terminal Illness, Continuity During Disability, Conversion On Withdrawal, and Conversion On Retirement. Applies per benefit.',
    data_source:
      'BaseRate * (1 + ContingencyLoading + VoluntaryLoading + SpecialLoadings)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'experience_adjusted_rate',
    data_type: 'number',
    data_description:
      'The Loaded Rate scaled by the scheme`s Experience Adjustment factor — the credibility-weighted blend of theoretical and observed claims experience. See gla_experience_adjustment and blended_gla_rate for the GLA-specific instance of the blend. Applies per benefit.',
    data_source: 'LoadedRate * ExperienceAdjustment',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'risk_premium',
    data_type: 'number',
    data_description:
      'The premium driven purely by risk: the experience-adjusted rate applied to the member`s sum assured (or annual income, for income-based benefits like TTD and PHI). Pre-loadings, pre-fees, pre-commission, pre-discount. Applies per benefit.',
    data_source: 'ExperienceAdjustedRate * SumAssured',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'basic_premium',
    data_type: 'number',
    data_description:
      'The residual portion of the final premium after the components that are presented explicitly on the premium breakdown (e.g. tax saver, binder, outsourcing, commission) have been separated out. It bundles together everything not surfaced on its own line — risk premium, expense loadings, profit margin, and any other components not shown separately. All explicit breakdown lines plus basic premium sum to the final premium.',
    data_source: 'FinalPremium - sum(explicit breakdown lines)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'office_premium',
    data_type: 'number',
    data_description:
      'The pre-discount, pre-commission premium charged to the scheme — the Risk Premium grossed up for the scheme-level loadings (expense, profit, admin, other, binder, outsourcing). Applies per benefit; rolls up to a scheme total.',
    data_source: 'RiskPremium / (1 - SchemeLoading)',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'discounted_rate',
    data_type: 'number',
    data_description:
      'The post-discount office premium expressed as a rate. Identical to office_premium when no discount is applied; otherwise computed under the discount method selected on the quote.',
    data_source:
      'Two methods, depending on the discount method chosen on the quote:\n\n  (1) Default (gross-up):\n      RiskPremium / (1 - (SchemeLoading + Discount))\n      — discount is folded into the loading denominator.\n\n  (2) Prorata:\n      RiskPremium * (1 - Discount) / (1 - SchemeLoading)\n      — discount is applied as a multiplicative factor on top of the standard office premium.\n\nBoth methods give the same result when the discount is zero.',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'final_rate',
    data_type: 'number',
    data_description:
      'The full rate actually charged to the scheme: the office premium (discounted if a discount applies, otherwise standard) plus commission. Commission is calculated progressively from the scheme-wide office premium before commission is added, then allocated back across members — it is not folded into the gross-up denominator.',
    data_source: 'OfficePremium (with or without discount) + Commission',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'region_loading',
    data_type: 'number',
    data_description:
      'Per-member rate adjustment reflecting geographic risk, resolved per member by region. Applied at the Base Rate stage. Applies per benefit.',
    data_source: 'Region Loading Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'industry_loading',
    data_type: 'number',
    data_description:
      'Per-member rate adjustment reflecting occupational / industry risk, resolved per member by industry. Applied at the Base Rate stage. Applies per benefit.',
    data_source: 'Industry Loading Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'contingency_loading',
    data_type: 'number',
    data_description:
      'Margin for adverse experience and parameter uncertainty. Applied at the Loaded Rate stage. Applies per benefit.',
    data_source: 'General Loadings Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'voluntary_loading',
    data_type: 'number',
    data_description:
      'Anti-selection loading applied where members elect voluntary cover. Applied at the Loaded Rate stage. Applies per benefit.',
    data_source: 'General Loadings Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'continuity_during_disability_loading',
    data_type: 'number',
    data_description:
      'Loading covering continued GLA cover while a member is disabled. Applied at the Loaded Rate stage.',
    data_source: 'General Loadings Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'conversion_on_withdrawal_loading',
    data_type: 'number',
    data_description:
      'Loading for the option to convert group cover to individual cover when a member withdraws from the scheme. Applies to GLA, PTD, and CI.',
    data_source: 'General Loadings Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'conversion_on_retirement_loading',
    data_type: 'number',
    data_description:
      'Loading for the option to convert GLA cover to individual cover at retirement.',
    data_source: 'General Loadings Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'expense_loading',
    data_type: 'number',
    data_description:
      'Insurer`s operating expenses expressed as a fraction of office premium. Component of scheme_loading.',
    data_source: 'Premium Loading Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'profit_loading',
    data_type: 'number',
    data_description:
      'Insurer`s required profit margin expressed as a fraction of office premium. Component of scheme_loading.',
    data_source: 'Premium Loading Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'admin_loading',
    data_type: 'number',
    data_description:
      'Scheme administration cost as a fraction of office premium. Component of scheme_loading.',
    data_source: 'Premium Loading Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'other_loading',
    data_type: 'number',
    data_description:
      'Catch-all for any additional scheme-level loadings not covered by the named components. Component of scheme_loading.',
    data_source: 'Premium Loading Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'scheme_loading',
    data_type: 'number',
    data_description:
      'Total scheme-level loading used to gross up Risk Premium into Office Premium. Commission is not included — it is added on top of the office premium separately.',
    data_source:
      'ExpenseLoading + ProfitLoading + AdminLoading + OtherLoading + BinderFee + OutsourcingFee',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'scheme_loading_after_discount',
    data_type: 'number',
    data_description:
      'The scheme loading adjusted for the discount applied to the quote. Discount is stored as a negative fraction (e.g. -0.05 for a 5% discount), so adding it shrinks the loading and therefore reduces the Discounted Rate.',
    data_source: 'SchemeLoading + Discount',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'binder_fee',
    data_type: 'number',
    data_description:
      'Fee paid to the binder holder for managing the scheme on behalf of the insurer. Included in the scheme loading denominator.',
    data_source: 'Binder Fees Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'outsourcing_fee',
    data_type: 'number',
    data_description:
      'Fee paid to an outsourced claims / administration provider. Included in the scheme loading denominator.',
    data_source: 'Binder Fees Table',
    data_source_type: 'User Input'
  },
  {
    data_variable: 'commission',
    data_type: 'number',
    data_description:
      'The amount payable to the intermediary (broker / agent). Calculated scheme-wide first from the office premium before commission is added, then allocated back to members so the gross-up is fair. Not folded into the scheme-loading denominator.',
    data_source:
      'Intermediary commission, applied progressively on top of the office premium based on the scheme-wide office premium before commission is added, then allocated back across members. Rate driven by the commission structure tiers configured on the quote.',
    data_source_type: 'Calculation Engine'
  },
  {
    data_variable: 'tax_saver',
    data_type: 'number',
    data_description:
      'An optional rider that adds a tax-adjusted retirement benefit on top of GLA. The loading is folded into the GLA Loaded Rate, but the resulting premium component is broken out separately on the breakdown so the GLA and Tax Saver lines do not double-count.',
    data_source:
      'Optional GLA rider for a tax-adjusted retirement benefit. Loaded into the GLA Loaded Rate via the Tax Saver loading; reported separately on the quote and not double-counted in GLA totals.',
    data_source_type: 'Calculation Engine'
  }
]
