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
    data_variable: 'free_cover_limit',
    data_type: 'number',
    data_description:
      'Similar to the MemberDistributionFreeCoverLimit, this limit is calculated based on the statistical distribution of members sums assured within the scheme. It is subject to the lesser of the scheme’s overall free cover limit or a predefined floor, which may be applied to maintain pricing competitiveness and manage underwriting risk.',
    data_source: 'Max (MemberDistributionFreeCoverLimit, MaxFreeCoverLimit )',
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
    data_variable: 'educator_risk_premium',
    data_type: 'number',
    data_description: 'Educator Risk Premium',
    data_source:
      '( Grade0SumAssured*Grade0RiskRate + Grade17SumAssured*Grade17RiskRate\n\t\t + Grade812SumAssured*Grade812RiskRate + TertiarySumAssured*TertiaryRiskRate ) *\n\t\t(ExpAdjLoadedGlaRate + ExpAdjLoadedTtdRate)',
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
  }
]
