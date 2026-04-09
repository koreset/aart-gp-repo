# Group Risk Pricing — Configuration and Calculation Reference

**Version:** v5.5.0
**Date:** 2026-03-30
**Audience:** Actuaries, product configurators, underwriters, and data analysts

This document is the field-level reference for setting up and operating the Group Risk Pricing
module. It covers configuration prerequisites, rate table mechanics, calculation formulas, and
output field semantics that are not described elsewhere.

For code-level production gaps and the high-level quote calculation flow, see
[gap_analysis.md](gap_analysis.md).

---

## Table of Contents

1. [System Configuration Prerequisites](#1-system-configuration-prerequisites)
2. [Rate Table Configuration](#2-rate-table-configuration)
3. [Income Benefit Configuration](#3-income-benefit-configuration)
4. [Quote-Level Settings](#4-quote-level-settings)
5. [Calculation Reference](#5-calculation-reference)

---

## 1. System Configuration Prerequisites

### 1.1 `risk_rate_code` — The Master Configuration Key

Every core lookup table in the pricing engine is keyed on `risk_rate_code`. This value on a
`GroupPricingQuote` record must match the `risk_rate_code` column in **all** of the tables
listed below before a quote can produce correct results.

`risk_rate_code` represents a product/jurisdiction/entity combination — for example, a South
African group risk product (`ZA_GRP_STD`) or a Namibian variant (`NA_GRP_STD`). Each distinct
product configuration requires its own complete set of seeded tables.

**Consequence of a missing entry:** rate lookups silently return 0 or an empty result — no
error is raised. A misconfigured `risk_rate_code` produces zero premiums with no user-visible
warning.

### 1.2 Required Tables Per `risk_rate_code`

The following tables **must** be seeded for a `risk_rate_code` before any quote can be priced.
Tables marked optional have defined fallback behaviour described in the relevant section.

| Table | Required? | Purpose |
|---|---|---|
| `group_pricing_parameters` | **Required** | All loadings, thresholds, FCL parameters, spouse age gap |
| `income_levels` | **Required** | Annual salary bands → integer level code used in every rate lookup |
| `restrictions` | **Required** | Maximum monthly benefit caps and maximum cover ages |
| `gla_rates` | Required if GLA benefit enabled | Base rates by ANB, gender, income level, waiting period, benefit type |
| `ptd_rates` | Required if PTD benefit enabled | Base rates by ANB, gender, income level, deferred period, disability definition |
| `ttd_rates` | Required if TTD benefit enabled | Base rates by ANB, gender, income level, waiting period, risk type, disability definition |
| `phi_rates` | Required if PHI benefit enabled | Base rates by ANB, gender, income level, waiting/deferred period, disability definition, escalation option, NRA |
| `ci_rates` | Required if CI benefit enabled | Base rates by ANB, gender, income level, benefit definition |
| `occupation_classes` | **Required** | Maps industry + scheme category → integer occupation class used in industry loading lookup |
| `group_pricing_age_bands` | **Required** | Age band labels (e.g., "25-29") applied to result records for reporting |
| `funeral_parameters` | Required if funeral benefit enabled | Average dependent/child age and count by member ANB |
| `tax_tables` | Optional | Marginal income tax bands. Falls back to gross salary if absent — see §3.1 |
| `tiered_income_replacements` | Conditional | Required when `phi_use_tiered_income_replacement_ratio = true` or `ttd_use_tiered_income_replacement_ratio = true` — see §3.2 |
| `reinsurance_structures` | Required if reinsurance configured | Quota share tier bounds and ceded proportions |
| `educator_benefit_structures` | Conditional | Required when `gla_educator_benefit = "Yes"` or `ptd_educator_benefit = "Yes"` |

### 1.3 `basis` — The Secondary Configuration Key

`GroupPricingParameters` is keyed on **both** `risk_rate_code` and `basis`. The `basis` value
on the quote request must match the `basis` column in `group_pricing_parameters`. Typical values
are `"standard"`, `"conservative"`, `"aggressive"` — the valid set depends on what has been
seeded and is not enforced by the engine.

If no `GroupPricingParameters` row matches the quote's `(basis, risk_rate_code)` pair, the
calculation returns an error. This is the only missing-configuration case that is explicitly
surfaced.

---

## 2. Rate Table Configuration

### 2.1 Income Levels

Every base rate table (GLA, PTD, CI, TTD, PHI, spouse GLA) includes an `income_level` lookup
key. The `income_levels` table maps each member's annual salary to that integer.

**`income_levels` table structure:**

| Column | Type | Description |
|---|---|---|
| `risk_rate_code` | string | Must match the quote |
| `year` | int | Rate year — must match the year of the quote's `basis` row |
| `min_income` | float | Lower bound of salary band (inclusive) |
| `max_income` | float | Upper bound of salary band (inclusive) |
| `level` | int | Integer code passed to rate table lookups |

**Lookup logic:** `GetIncomeLevel` iterates the income level rows for the quote's `risk_rate_code`
and returns the first `level` where `min_income ≤ annualSalary ≤ max_income`. If no row matches,
it returns `0`.

**Silent failure:** If `level = 0` and no rate table row exists for `income_level = 0`, the base
rate for every benefit returns 0. No error is raised. Ensure the income bands cover the full
salary range of the scheme's membership and that the integer `level` values used in `income_levels`
exactly match the `income_level` column in the individual rate tables.

**Example:**

| min_income | max_income | level |
|---|---|---|
| 0 | 250 000 | 1 |
| 250 001 | 500 000 | 2 |
| 500 001 | 750 000 | 3 |
| 750 001 | 0 | 4 |

A member with annual salary R420 000 maps to level 2. The engine then looks up
`gla_rates` where `income_level = 2`.

### 2.2 `GroupPricingParameters` — Full Field Reference

`GroupPricingParameters` is the primary configuration record for a `(basis, risk_rate_code)` pair.

| Field | Type | Description |
|---|---|---|
| `basis` | string | Pricing basis — must match `GroupPricingQuote.basis` |
| `risk_rate_code` | string | Must match `GroupPricingQuote.risk_rate_code` |
| `treaty_code` | string | Links to the reinsurance treaty for this basis/code pair; used when loading `GroupPricingReinsuranceStructure` |
| `educator_benefit_code` | string | Selects the `educator_benefit_structures` row; only loaded when `gla_educator_benefit = "Yes"` or `ptd_educator_benefit = "Yes"` |
| `spouse_age_gap` | int | Age gap between member and spouse. Female member → spouse ANB = member ANB + gap. Male member → spouse ANB = member ANB − gap |
| `min_age` / `max_age` | int | Age clamps applied when deriving spouse ANB |
| `full_credibility_threshold` | float | Minimum weighted life-years required for 100% credibility weight. Set to 0 to disable experience rating entirely |
| `free_cover_limit_scaling_factor` | float | Multiplier in FCL formula: `scalingFactor × sqrt(n) × meanSalary` — see §4.2 |
| `free_cover_limit_percentile` | float | Salary percentile cap on FCL (e.g., 0.95 = 95th percentile). Set to 0 to disable |
| `free_cover_limit_nearest_multiple` | float | Round calculated FCL up to nearest multiple (e.g., 50 000). Set to 0 to disable |
| `global_gla_experience_rate` | float | Fallback crude experience rate for GLA when no `GroupPricingClaimsExperience` records exist |
| `global_ptd_experience_rate` | float | Same for PTD |
| `global_ci_experience_rate` | float | Same for CI |
| `is_lumpsum_reins_gla_dependent` | bool | If `true`, PTD / CI / SGLA reinsurance cession uses the GLA retained proportion instead of their own tier levels |
| `gla_terminal_illness_loading_rate` | float | Additive loading on GLA rate when `SchemeCategory.gla_terminal_illness_benefit = "Yes"` |
| `ptd_accelerated_benefit_discount` | float | Proportional discount on PTD rate when `ptd_benefit_type = "Accelerated"`. Applied as `rate × (1 − discount)` |
| `ci_accelerated_benefit_discount` | float | Same for CI accelerated benefit |
| `ttd_number_monthly_payments` | float | Number of monthly TTD payments per claim year; multiplied against TTD risk premium |
| `maximum_contribution_waiver` | float | Monthly cap on PHI pension contribution waiver at the parameters level |
| `medical_aid_waiver_proportion` | float | Proportion of monthly salary used to estimate medical aid premium (for PHI medical aid waiver calculation) |
| `medical_aid_waiver_amount` | float | Fixed rand amount added to the proportional estimate for medical aid waiver |
| `maximum_medical_aid_waiver` | float | Monthly cap on PHI medical aid waiver |
| `minimum_profitability_loading` | float | Floor for total loading. `TotalLoading = max(sum_of_all_loadings, minimum_profitability_loading)` — see §5.2 |
| `expense_loading` | float | As a decimal (e.g., 0.08 = 8%). Stored in parameters; multiplied × 100 before display |
| `admin_loading` | float | Separate from expense loading — covers in-scheme administration costs |
| `commission_loading` | float | Set to 0 automatically when `distribution_channel = "direct"` |
| `contingency_loading` | float | Stored in parameters but **currently excluded** from the TotalLoading sum — see §5.2 |
| `profit_loading` | float | Insurer profit margin loading |
| `other_loading` | float | Miscellaneous loading |
| `discount` | float | Enter as a positive value. The engine subtracts it from the loading sum — see §5.2 |
| `premium_rates_guaranteed_period_months` | int | Stored for reference; not currently applied in any premium calculation |
| `quote_validity_period_months` | int | Stored for reference; not currently applied in any premium calculation |
| `risk_margin_rate` | float | Stored for reference; not currently applied in the main member-level calculation |
| `reinsurer_profit_loading` | float | Stored for reference; not currently applied in the main member-level calculation |
| `annual_expense_amount` | float | Fixed annual expense amount (rand value). Stored for reference; not currently applied in member-level premium. Distinct from `expense_loading` (percentage) |

### 2.3 Rate Table Lookup Keys

Each benefit's base rate is looked up by combining several dimensions. The string values used
in `SchemeCategory` for the fields below must **exactly** match the corresponding column values
in the rate table — including case. A mismatch produces a zero base rate without an error.

| Benefit | SchemeCategory field | Rate table column | Example values |
|---|---|---|---|
| GLA | `gla_benefit_type` | `benefit_type` | Depends on `risk_rate_code` seeding |
| GLA | `gla_waiting_period` | `waiting_period` | `0`, `3`, `6` (months) |
| PTD | `ptd_risk_type` | `risk_type` | Depends on `risk_rate_code` seeding |
| PTD | `ptd_deferred_period` | `deferred_period` | `0`, `3`, `6` (months) |
| PTD | `ptd_disability_definition` | `disability_definition` | Depends on `risk_rate_code` seeding |
| CI | `ci_benefit_definition` | `benefit_definition` | Depends on `risk_rate_code` seeding |
| TTD | `ttd_risk_type` | `risk_type` | Depends on `risk_rate_code` seeding |
| TTD | `ttd_waiting_period` | `waiting_period` | `0`, `7`, `14` (days) |
| TTD | `ttd_disability_definition` | `disability_definition` | Depends on `risk_rate_code` seeding |
| PHI | `phi_risk_type` | `risk_type` | Depends on `risk_rate_code` seeding |
| PHI | `phi_waiting_period` | `waiting_period` | `4`, `13`, `26` (weeks) |
| PHI | `phi_deferred_period` | `deferred_period` | `0`, `3`, `6` (months) |
| PHI | `phi_disability_definition` | `disability_definition` | Depends on `risk_rate_code` seeding |
| PHI | `phi_benefit_escalation` | `benefit_escalation_option` | Depends on `risk_rate_code` seeding |
| PHI | `phi_normal_retirement_age` | `normal_retirement_age` | `60`, `63`, `65` |

**Note on `phi_normal_retirement_age`:** This is a PHI-specific NRA used only for the PHI rate
table lookup. It is separate from `GroupPricingQuote.normal_retirement_age`, which is used to
set the `ExceedsNormalRetirementAgeIndicator` flag on member results.

---

## 3. Income Benefit Configuration

### 3.1 Tax Tables and Take-Home Pay

The PHI and TTD monthly income benefit is capped at the member's take-home pay (net of income
tax). The engine computes take-home pay by applying marginal tax bands from the `tax_tables`
table before performing the income replacement calculation.

**`tax_tables` structure:**

| Column | Type | Description |
|---|---|---|
| `risk_rate_code` | string | Must match the quote |
| `level` | int | Band order — rows must be processed in ascending level order |
| `min` | float | Lower bound of this salary band |
| `max` | float | Upper bound. Set to `0` to indicate an open-ended top band |
| `tax_rate` | float | Marginal rate applied to the salary portion within this band (e.g., 0.18 = 18%) |

**How the calculation works:**

The engine iterates bands in ascending `level` order. For each band, it taxes only the portion
of the member's remaining unallocated salary that falls within that band's `[min, max)` range.
An open-ended top band (`max = 0`) absorbs all remaining salary above `min`.

```
takeHomePay = annualSalary − Σ (portionInBand_i × taxRate_i)
monthlyTakeHomePay = takeHomePay / 12
```

**Fallback:** If no rows exist in `tax_tables` for the quote's `risk_rate_code`, take-home pay
defaults to gross annual salary (effective tax rate = 0%). The income benefit then equals
the flat percentage of gross salary — identical to the pre-tax-table behaviour.

**Example — South African tax year 2025/26:**

| level | min | max | tax_rate |
|---|---|---|---|
| 1 | 0 | 95 750 | 0.00 |
| 2 | 95 750 | 365 000 | 0.18 |
| 3 | 365 000 | 550 000 | 0.26 |
| 4 | 550 000 | 0 | 0.36 |

Member with annual salary R420 000:
- Band 1: R95 750 × 0% = R0
- Band 2: R269 250 × 18% = R48 465
- Band 3: R55 000 × 26% = R14 300
- Total tax = R62 765 | Take-home = R357 235 | Monthly = R29 770

### 3.2 Income Replacement: Flat Rate vs Tiered Schedule

Two boolean flags on `SchemeCategory` control whether PHI and TTD income replacement uses a
flat percentage or a tiered salary-band schedule.

| Field | Type | Default | Effect |
|---|---|---|---|
| `phi_use_tiered_income_replacement_ratio` | bool | false | Use tiered table for PHI income replacement |
| `ttd_use_tiered_income_replacement_ratio` | bool | false | Use tiered table for TTD income replacement |

#### Flat Rate (flag = `false`)

```
annualReplacementIncome = annualSalary × (incomeReplacementPercentage / 100)
unScaledMonthlyIncome   = min(monthlyTakeHomePay, annualReplacementIncome / 12)
```

`incomeReplacementPercentage` comes from `SchemeCategory.phi_income_replacement_percentage` or
`ttd_income_replacement_percentage` (entered as a percentage, e.g., 75 for 75%).

#### Tiered Replacement (flag = `true`)

```
annualReplacementIncome = Σ portionOfSalaryInTier_i × replacementRatio_i
unScaledMonthlyIncome   = min(monthlyTakeHomePay, annualReplacementIncome / 12)
```

The tiered schedule is stored in `tiered_income_replacements`:

| Column | Type | Description |
|---|---|---|
| `risk_rate_code` | string | Must match the quote |
| `annual_lower_bound` | float | Lower salary bound for this tier (inclusive) |
| `annual_upper_bound` | float | Upper salary bound for this tier |
| `income_replacement_ratio` | float | Decimal replacement ratio for this tier (e.g., 0.75 = 75%) |

The engine applies the same band-iteration logic as the tax table: only the portion of salary
within each tier's bounds is multiplied by that tier's ratio.

**Example:**

| annual_lower_bound | annual_upper_bound | income_replacement_ratio |
|---|---|---|
| 0 | 500 000 | 0.75 |
| 500 000 | 1 000 000 | 0.50 |
| 1 000 000 | 0 | 0.25 |

Member with annual salary R700 000:
- Tier 1: R500 000 × 75% = R375 000
- Tier 2: R200 000 × 50% = R100 000
- Annual replacement income = R475 000

**Silent failure:** If `phi_use_tiered_income_replacement_ratio = true` but no rows exist in
`tiered_income_replacements` for the quote's `risk_rate_code`, `annualReplacementIncome`
computes to R0, making `PhiCappedIncome = 0`. No error is raised. Always seed the tiered
table before enabling the flag.

### 3.3 Restriction Table — Maximum Caps and Age Ceilings

The `restrictions` table defines hard caps applied **after** the income replacement calculation.
These are separate from the salary-multiple and income replacement logic — they represent absolute
regulatory or product limits.

**Monthly income caps** — applied as `CappedIncome = min(unScaledIncome, cap)`:

| Field | Applies to |
|---|---|
| `phi_maximum_monthly_benefit` | PHI `PhiCappedIncome` |
| `ttd_maximum_monthly_benefit` | TTD `TtdCappedIncome` |
| `phi_maximum_monthly_contribution_waiver` | PHI per-member pension contribution waiver |
| `max_medical_aid_waiver` | PHI medical aid premium waiver |
| `spouse_gla_maximum_benefit` | Spouse GLA sum assured |
| `severe_illness_maximum_benefit` | CI sum assured |

**Maximum cover ages** — members whose age-next-birthday (ANB) exceeds these values are
excluded from rating for the corresponding benefit:

| Field | Benefit |
|---|---|
| `gla_max_cover_age` | GLA |
| `ptd_max_cover_age` | PTD |
| `ci_max_cover_age` | CI |
| `phi_max_cover_age` | PHI |
| `ttd_max_cover_age` | TTD |
| `fun_max_cover_age` | Family funeral (all variants) |

**Note on `NormalRetirementAge`:** `GroupPricingQuote.normal_retirement_age` sets the
`exceeds_normal_retirement_age_indicator` flag on member results but does **not** currently
exclude the member from rating. The restriction table's max cover ages are the operative
exclusion mechanism.

---

## 4. Quote-Level Settings

### 4.1 Distribution Channel and Commission

`GroupPricingQuote.distribution_channel` controls how commission is applied to the quote.

| Value | Commission treatment |
|---|---|
| `broker` | `CommissionLoading` from `GroupPricingParameters` applied in full |
| `direct` | `CommissionLoading` forced to 0 regardless of the parameters value |
| `binder` | Commission applied as per parameters (same as `broker`) |
| `tied_agent` | Commission applied as per parameters (same as `broker`) |

Selecting `direct` for a scheme where the parameters record has a non-zero `commission_loading`
will silently reduce the quoted premium. Verify the channel before running a calculation.

### 4.2 Free Cover Limit Calculation

The free cover limit (FCL) is the maximum sum assured that can be provided without individual
evidence of insurability. The engine calculates it statistically from the membership salary
distribution unless overridden at the quote level.

**Formula:**

```
scaledFCL  = free_cover_limit_scaling_factor × sqrt(memberCount) × meanAnnualSalary

if free_cover_limit_percentile > 0:
    nthPercentile = Nth percentile of the member salary distribution
    calculatedFCL = min(scaledFCL, nthPercentile)
else:
    calculatedFCL = scaledFCL

if free_cover_limit_nearest_multiple > 0:
    calculatedFCL = ceil(calculatedFCL / nearest_multiple) × nearest_multiple

if GroupPricingQuote.free_cover_limit > 0:
    appliedFCL = min(calculatedFCL, GroupPricingQuote.free_cover_limit)
else:
    appliedFCL = calculatedFCL
```

**Parameter guide:**

| Parameter | Typical value | Set to 0 to… |
|---|---|---|
| `free_cover_limit_scaling_factor` | 1.5 – 3.0 | Disable the statistical formula (FCL becomes 0) |
| `free_cover_limit_percentile` | 0.90 – 0.95 | Remove the percentile cap — FCL is unconstrained by distribution |
| `free_cover_limit_nearest_multiple` | 10 000 – 100 000 | Disable rounding |

If a non-zero `GroupPricingQuote.free_cover_limit` is provided, it acts as a ceiling on the
calculated value. The lower of the calculated and the specified FCL is applied.

### 4.3 Indicative Member Data Mode

When `GroupPricingQuote.member_indicative_data = true`, the engine uses a single representative
data set instead of individual member records. This is intended for early-stage quotes where a
full member census is not yet available.

**Key differences from full member mode:**

| Behaviour | Full member mode | Indicative mode |
|---|---|---|
| Member data source | `gpricing_member_datas` table (one row per member) | Single `MemberIndicativeDataSet` record |
| `indicativeRatesCount` | 1 (each result row = one member) | `MemberDataCount` (each result row = entire portfolio) |
| Output field meaning | Individual member values | Portfolio totals (all members rated as one representative) |
| FCL calculation | Statistical formula using salary distribution | `annualSalary × GlaSalaryMultiple` (no distribution available) |
| Age used | Calculated from each member's date of birth | `MemberIndicativeDataSet.MemberAverageAge` |

**Important:** The output fields `GlaSumAssured`, `PhiIncome`, etc. are multiplied by
`indicativeRatesCount` in indicative mode. A result showing `GlaSumAssured = 50 000 000` in
indicative mode means the portfolio total, not an individual member's sum assured.

Switching a quote from indicative to full member mode will change result magnitudes structurally.
Do not mix indicative and full member outputs when comparing two runs of the same quote.

### 4.4 `SchemeCategory` — Alias Fields

Seven alias fields on `SchemeCategory` set the display names for benefits in the quote PDF
and any client-facing report. They have no effect on the rating calculation.

| Field | Default display label replaced |
|---|---|
| `gla_alias` | "GLA" / "Group Life Assurance" |
| `sgla_alias` | "SGLA" / "Spouse GLA" |
| `ptd_alias` | "PTD" / "Permanent Total Disability" |
| `ci_alias` | "CI" / "Critical Illness" |
| `phi_alias` | "PHI" / "Permanent Health Insurance" |
| `ttd_alias` | "TTD" / "Temporary Total Disability" |
| `family_funeral_alias` | "Family Funeral" |

### 4.5 `SchemeCategory` — PHI and TTD Additional Fields

| Field | Description |
|---|---|
| `phi_number_monthly_payments` | Maximum number of PHI monthly income payments per claim. Stored on the category; currently informational only — the active limit is in `GroupPricingParameters.ttd_number_monthly_payments` for TTD |
| `phi_max_premium_waiver` | Category-level monthly cap on the PHI pension contribution waiver. Applied in addition to the parameters-level `maximum_contribution_waiver` |
| `phi_medical_aid_premium_waiver` | Set to `"Yes"` to include a medical aid premium waiver in the PHI monthly benefit. The waiver amount is calculated from `medical_aid_waiver_proportion × salary + medical_aid_waiver_amount`, capped at `restrictions.max_medical_aid_waiver` |
| `phi_premium_waiver` | Set to `"Yes"` to include a pension contribution waiver in the PHI monthly benefit. The waiver amount is `member.contribution_waiver_proportion × monthly_salary`, capped at `GroupPricingParameters.maximum_contribution_waiver` |

### 4.6 `use_global_salary_multiple`

When `GroupPricingQuote.use_global_salary_multiple = true`:
- All members use `SchemeCategory.gla_salary_multiple`, `ptd_salary_multiple`, `sgla_salary_multiple`,
  and `ci_critical_illness_salary_multiple`.
- The `benefits.*_multiple` fields on individual `GPricingMemberData` records are **ignored entirely**.

When `false`, per-member multiples from `GPricingMemberData.benefits` are used instead.

---

## 5. Calculation Reference

### 5.1 Per-Member Data — Behavioural Notes

| Field | Behavioural note |
|---|---|
| `contribution_waiver_proportion` | Represents the member's pension fund contribution rate as a decimal (e.g., 0.075 = 7.5% of salary). Applied as a monthly waiver when `SchemeCategory.phi_premium_waiver = "Yes"`. This must be set per member in the uploaded CSV — it defaults to 0 if omitted |
| `occupational_class` | Stored on the member record but **not used** in rate table lookups. Occupation class for rating is derived at the quote level from `GroupPricingQuote.industry` via the `occupation_classes` table. Per-member occupational class has no current effect on pricing |
| `entry_date` / `exit_date` | Used for experience rating period calculations. `exit_date` defaults to `1900-01-01` when not supplied — this sentinel value is a known issue (see GAP-07 in gap_analysis.md). Do not rely on exit date calculations until this is resolved |
| `is_original_member` | `true` for members present at scheme inception; `false` for mid-term joiners. Used to distinguish member cohorts in renewal quotes but does not currently affect the base rate |

### 5.2 Loading Formula and Office Premium

All premiums are converted from risk premiums (pure cost of claims) to office premiums (amount
charged to the scheme) by dividing by `(1 − TotalLoading)`.

**TotalLoading formula:**

```
TotalLoading = max(
    ExpenseLoading + AdminLoading + CommissionLoading + ProfitLoading + OtherLoading − Discount,
    MinimumProfitabilityLoading
)

OfficePremium = RiskPremium / (1 − TotalLoading)
```

**Important notes on individual loading components:**

| Component | Note |
|---|---|
| `expense_loading` | Covers claims processing and general overhead |
| `admin_loading` | In-scheme administration costs — distinct from `expense_loading`. Both are included in `TotalLoading` |
| `commission_loading` | Set to 0 for `distribution_channel = "direct"` regardless of the parameters value — see §4.1 |
| `discount` | Enter as a **positive** value in `GroupPricingParameters`. The formula subtracts it, reducing the total loading. A value of 0.02 = 2% discount |
| `contingency_loading` | Stored in `GroupPricingParameters` but **currently excluded** from the `TotalLoading` sum. It does not affect quoted premiums in the current version |
| `minimum_profitability_loading` | Ensures `TotalLoading` never falls below this floor, regardless of the sum of individual components |

**Warning:** If `TotalLoading ≥ 1.0`, the `OfficePremium` formula produces a zero or negative
result. The engine does not guard against this. Verify that the sum of all loadings remains below
1.0 when configuring `GroupPricingParameters`.

### 5.3 Experience Rating — Data Requirements

Experience rating blends the theoretical (table) rate with an empirical crude rate derived from
the scheme's own claims history. It activates only when all of the following conditions are met:

1. `GroupPricingQuote.experience_rating = "Yes"`
2. At least one `GroupPricingClaimsExperience` record exists for the quote
3. `GroupPricingParameters.full_credibility_threshold > 0`

**Required fields on `GroupPricingClaimsExperience`:**

| Field | Description |
|---|---|
| `number_of_members` | Number of exposed lives in the experience period |
| `claim_amount` | Total claims paid in the period |
| `time_period_years` | Length of observation period in years |
| `weighting` | Credibility weight for this period (e.g., 1.0 for all periods equal) |

**Credibility formula:**

```
weightedLifeYears = Σ (number_of_members_i × time_period_years_i × weighting_i)
credibility       = min( sqrt(weightedLifeYears / full_credibility_threshold), 1.0 )
adjustedRate      = theoreticalRate × (1 − credibility) + credibility × experienceCrudeRate
```

**Manual override:** If a non-zero `credibility` value is passed directly to the
`CalculateGroupPricingQuote` API call, it replaces the calculated credibility entirely.
The `manually_added_credibility` field on `MemberRatingResultSummary` stores this override.

**PHI and TTD:** Experience rating applies only to **GLA, PTD, and CI**. PHI and TTD
experience adjustments are always set to 1.0 (no blend) — there is no TTD or PHI crude
rate calculation regardless of claims data supplied.

**Global fallback rates:** When no `GroupPricingClaimsExperience` records exist (or
`experience_rating = "No"`), `global_gla_experience_rate`, `global_ptd_experience_rate`,
and `global_ci_experience_rate` from `GroupPricingParameters` are used as the crude rates
in the blending formula. Set these to the same value as the typical market crude rate for the
scheme size to produce a sensible non-experience-rated result.

### 5.4 Understanding the Output Fields (`MemberRatingResult`)

#### Uncapped vs Capped Fields

For sum assured benefits, two variants of each field are produced:

| Uncapped field | Capped field | Capped by |
|---|---|---|
| `gla_sum_assured` | `gla_capped_sum_assured` | `AppliedFreeCoverLimit` |
| `ptd_sum_assured` | `ptd_capped_sum_assured` | GLA capped SA (and FCL for standalone PTD) |
| `ci_sum_assured` | `ci_capped_sum_assured` | `restrictions.severe_illness_maximum_benefit` |
| `spouse_gla_sum_assured` | `spouse_gla_capped_sum_assured` | `restrictions.spouse_gla_maximum_benefit` |

For income benefits:

| Uncapped field | Capped field | Capped by |
|---|---|---|
| `phi_income` | `phi_capped_income` | `restrictions.phi_maximum_monthly_benefit` |
| `ttd_income` | `ttd_capped_income` | `restrictions.ttd_maximum_monthly_benefit` |

The **capped** fields are used for all premium calculations. If `phi_income ≠ phi_capped_income`,
the restriction cap is binding for that member — the PHI benefit is lower than the replacement
income formula would otherwise produce.

#### Theoretical vs Experience-Adjusted Fields

All premium fields exist in two variants:

| Theoretical (table rate) | Experience-adjusted |
|---|---|
| `gla_risk_premium` | `exp_adj_gla_risk_premium` |
| `gla_office_premium` | `exp_adj_gla_office_premium` |
| `ttd_risk_premium` | `exp_adj_ttd_risk_premium` |
| *(all benefits follow this pattern)* | |

The **experience-adjusted** fields incorporate the credibility blend. These are the premiums
charged to the scheme. The theoretical fields represent the pure rate-table cost before any
experience adjustment — useful for comparing the scheme's actual cost against the expected
market rate.

When `experience_rating = "No"` or `credibility = 0`, both variants are identical.

#### `income_replacement_ratio`

Stored as a **decimal** (0.0–1.0), not the percentage entered on `SchemeCategory`. It is
derived from `phi_income_replacement_percentage` only — TTD does not populate this field.
For example, a 75% replacement percentage produces `income_replacement_ratio = 0.75`.

#### Indicator Flags

| Field | Meaning |
|---|---|
| `exceeds_normal_retirement_age_indicator` | `1` if the member's ANB > `GroupPricingQuote.normal_retirement_age`. This is a **flag only** — the member is still rated. The restriction table's max cover ages are the operative exclusion |
| `exceeds_free_cover_limit_indicator` | `1` if `gla_sum_assured > GroupPricingQuote.free_cover_limit`. Flag only — does not change the premium calculation |

#### Educator Benefit Fields

The fields `grade_0_sum_assured`, `grade_1_7_sum_assured`, `grade_8_12_sum_assured`,
`tertiary_sum_assured` (and their corresponding rate and premium fields) are non-zero only
when `SchemeCategory.gla_educator_benefit = "Yes"` or `ptd_educator_benefit = "Yes"`. The
grades refer to South African schooling levels:

| Grade code | Description |
|---|---|
| Grade 0 | Pre-primary / ECD |
| Grade 1–7 | Primary school |
| Grade 8–12 | High school |
| Tertiary | University / college |

The sum assured for each grade is `max_tuition_per_year × max_coverage_years × (1 + book_allowance_proportion + accommodation_allowance_proportion)` from `EducatorBenefitStructure`.
