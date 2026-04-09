# Three-Tier Reinsurance Treaty Structure

## Overview

AART now supports sophisticated three-tier reinsurance structures that allow for:
- **Sum Assured Tiers**: Three levels with configurable bounds and cession proportions
- **Income Tiers**: Three levels based on annual income/salary with separate cession proportions
- **Multi-Reinsurer Support**: Support for lead reinsurer plus up to 3 non-lead reinsurers with different share percentages

## Database Schema Changes

### New Fields in `reinsurance_treaties` Table

#### Basic Treaty Configuration
- `treaty_code` (VARCHAR): Unique treaty code identifier
- `risk_premium_basis_indicator` (VARCHAR): `default`, `flat_rate`, or `custom`
- `flat_annual_reins_prem_rate` (DECIMAL): Applied when basis is `flat_rate`

#### Sum Assured Level 1
- `level1_lowerbound` (DECIMAL): e.g., 0
- `level1_upperbound` (DECIMAL): e.g., 500,000
- `level1_ceded_proportion` (DECIMAL): e.g., 20% (stored as 20.00)

#### Sum Assured Level 2
- `level2_lowerbound` (DECIMAL): e.g., 500,001
- `level2_upperbound` (DECIMAL): e.g., 1,500,000
- `level2_ceded_proportion` (DECIMAL): e.g., 40%

#### Sum Assured Level 3
- `level3_lowerbound` (DECIMAL): e.g., 1,500,001
- `level3_upperbound` (DECIMAL): e.g., 5,000,000
- `level3_ceded_proportion` (DECIMAL): e.g., 60%

#### Income Level 1
- `income_level1_lowerbound` (DECIMAL): e.g., 0
- `income_level1_upperbound` (DECIMAL): e.g., 200,000
- `income_level1_ceded_proportion` (DECIMAL): e.g., 15%

#### Income Level 2
- `income_level2_lowerbound` (DECIMAL): e.g., 200,001
- `income_level2_upperbound` (DECIMAL): e.g., 500,000
- `income_level2_ceded_proportion` (DECIMAL): e.g., 30%

#### Income Level 3
- `income_level3_lowerbound` (DECIMAL): e.g., 500,001
- `income_level3_upperbound` (DECIMAL): e.g., 2,000,000
- `income_level3_ceded_proportion` (DECIMAL): e.g., 50%

#### Multi-Reinsurer Structure
- `lead_reinsurer_code` (VARCHAR): Code for lead reinsurer
- `lead_reinsurer_share` (DECIMAL): Lead reinsurer's share % (e.g., 50%)
- `non_lead_reinsurer1_code` (VARCHAR): First follower reinsurer code
- `non_lead_reinsurer1_share` (DECIMAL): e.g., 25%
- `non_lead_reinsurer2_code` (VARCHAR): Second follower reinsurer code
- `non_lead_reinsurer2_share` (DECIMAL): e.g., 15%
- `non_lead_reinsurer3_code` (VARCHAR): Third follower reinsurer code
- `non_lead_reinsurer3_share` (DECIMAL): e.g., 10%
- `ceding_commission` (DECIMAL): Commission paid by reinsurer to cedant %

## Cession Calculation Logic

### Priority Order

The cession calculation engine in `backend/services/reinsurance_cession.go` follows this priority:

1. **Check for Three-Tier Structure**: If `level1_upperbound > 0` and `level1_ceded_proportion > 0`, use tiered cession
2. **Fallback to Legacy Logic**: Otherwise, use the original treaty type logic (proportional, surplus, XL, etc.)

### Tiered Cession Algorithm

```go
func calculateTieredCession(sumAssured float64, treaty models.ReinsuranceTreaty) CessionResult
```

**Logic**:
- Checks which tier the sum assured falls into (Level 3 → Level 2 → Level 1)
- Applies the corresponding `ceded_proportion` percentage
- Returns retention amount, ceded amount, cession percentage, and tier identifier

**Example**:
```
Sum Assured: R 750,000
Level 1: R 0 - R 500,000 @ 20% cession
Level 2: R 500,001 - R 1,500,000 @ 40% cession
Level 3: R 1,500,001 - R 5,000,000 @ 60% cession

Result: Falls in Level 2 → 40% cession
Ceded: R 300,000
Retained: R 450,000
```

### Income-Based Cession

```go
func CalculateMemberCessionWithIncome(annualIncome float64, treaty models.ReinsuranceTreaty) CessionResult
```

- Similar logic but uses income levels instead of sum assured levels
- Useful for salary-linked benefits where cession varies by income band
- Returns `cession_basis` as `income_level1`, `income_level2`, or `income_level3`

### Outside Range Handling

If sum assured or income falls outside all configured tiers:
- **Cession Basis**: `tiered_outside_range` or `income_outside_range`
- **Ceded Amount**: 0
- **Retention**: Full amount

## API Changes

### Create Treaty Request

**POST** `/group-pricing/reinsurance/treaties`

New fields accepted in request body (all optional):
```json
{
  "treaty_number": "TREATY-2026-001",
  "treaty_name": "Group Life RI Treaty 2026",
  "reinsurer_name": "Munich Re",
  "treaty_code": "GLRT2026",
  "risk_premium_basis_indicator": "default",
  "flat_annual_reins_prem_rate": 0,
  "level1_lowerbound": 0,
  "level1_upperbound": 500000,
  "level1_ceded_proportion": 20,
  "level2_lowerbound": 500001,
  "level2_upperbound": 1500000,
  "level2_ceded_proportion": 40,
  "level3_lowerbound": 1500001,
  "level3_upperbound": 5000000,
  "level3_ceded_proportion": 60,
  "income_level1_lowerbound": 0,
  "income_level1_upperbound": 200000,
  "income_level1_ceded_proportion": 15,
  "income_level2_lowerbound": 200001,
  "income_level2_upperbound": 500000,
  "income_level2_ceded_proportion": 30,
  "income_level3_lowerbound": 500001,
  "income_level3_upperbound": 2000000,
  "income_level3_ceded_proportion": 50,
  "lead_reinsurer_code": "MURE",
  "lead_reinsurer_share": 50,
  "non_lead_reinsurer1_code": "SWRE",
  "non_lead_reinsurer1_share": 30,
  "non_lead_reinsurer2_code": "HANN",
  "non_lead_reinsurer2_share": 20,
  "ceding_commission": 10
}
```

### Update Treaty Request

**PUT** `/group-pricing/reinsurance/treaties/:id`

All new fields can be updated independently. The update service uses `>= 0` checks to allow zero values.

## Frontend Changes

### RITreatyManagement.vue Form Updates

New sections added to the treaty form dialog:

1. **Three-Tier Reinsurance Structure**
   - Treaty code input
   - Risk premium basis selector
   - Flat annual rate (enabled when basis = `flat_rate`)

2. **Sum Assured Tiers (Levels 1-3)**
   - Each level has: Lower Bound, Upper Bound, Ceded Proportion %

3. **Income Tiers (Income Levels 1-3)**
   - Each level has: Lower Bound, Upper Bound, Ceded Proportion %

4. **Multi-Reinsurer Structure**
   - Lead reinsurer code and share %
   - Up to 3 non-lead reinsurers with codes and shares
   - Ceding commission %

All fields are optional and have proper validation.

## Bordereaux Generation Impact

### Member Census Bordereaux

When generating member bordereaux:
- If three-tier structure is configured, the engine uses `calculateTieredCession()`
- The `cession_basis` field in `ri_bordereaux_member_rows` will show:
  - `tiered_level1`, `tiered_level2`, or `tiered_level3` for sum assured tiers
  - `income_level1`, `income_level2`, or `income_level3` for income tiers
  - `tiered_outside_range` if no tier matches

This allows reinsurers to see which tier was applied for each member.

### Claims Bordereaux

Claims cession continues to use the existing `CalculateClaimCession()` logic:
- XL treaties: cession based on claim amount vs XL retention/limit
- Proportional: uses legacy `cession_percentage` or falls back to tiered if configured

## Migration Instructions

1. **Database Migration**:
   ```bash
   cd backend
   # Run the migration SQL file
   mysql -u root -p new_aart < migrations/add_three_tier_treaty_structure.sql
   # OR for PostgreSQL:
   psql -U postgres -d new_aart -f migrations/add_three_tier_treaty_structure.sql
   ```

2. **Backend Deployment**:
   ```bash
   cd backend
   go build -o aart_api
   ./aart_api
   ```

3. **Frontend Deployment**:
   ```bash
   cd frontend
   npm run build
   # Or for development:
   npm run dev
   ```

## Backward Compatibility

- **Existing Treaties**: All existing treaties continue to work with legacy `cession_percentage`
- **Migration**: No data migration needed - new fields default to 0
- **Cession Logic**: System automatically detects if three-tier structure is configured
  - If `level1_upperbound > 0` and `level1_ceded_proportion > 0` → Uses tiered logic
  - Otherwise → Uses legacy treaty type logic

## Validation Considerations

When using three-tier structures:

1. **Tier Continuity**: Ensure upper bound of Level N equals (or is close to) lower bound of Level N+1
2. **Cession Progression**: Typically, higher tiers have higher cession percentages
3. **Multi-Reinsurer Shares**: Total shares should add up to 100% (not enforced but recommended)
4. **Income vs Sum Assured**: Can configure one, both, or neither tier structure independently

## Example Scenarios

### Scenario 1: Simple Three-Tier Sum Assured

```
Treaty: Group Life Proportional RI
Level 1: R 0 - R 500,000 @ 20% cession
Level 2: R 500,001 - R 1,500,000 @ 40% cession
Level 3: R 1,500,001 - R 10,000,000 @ 60% cession

Member A: Sum Assured R 300,000 → Level 1 → 20% ceded = R 60,000
Member B: Sum Assured R 800,000 → Level 2 → 40% ceded = R 320,000
Member C: Sum Assured R 3,000,000 → Level 3 → 60% ceded = R 1,800,000
```

### Scenario 2: Multi-Reinsurer Panel

```
Treaty: Group Life RI Treaty 2026
Lead Reinsurer: Munich Re @ 50%
Follower 1: Swiss Re @ 30%
Follower 2: Hannover Re @ 20%

If total cession is R 1,000,000:
- Munich Re: R 500,000 (50%)
- Swiss Re: R 300,000 (30%)
- Hannover Re: R 200,000 (20%)
```

### Scenario 3: Income-Based Cession

```
Treaty: Salary-Linked Death Benefit RI
Income Level 1: R 0 - R 200,000 @ 10% cession
Income Level 2: R 200,001 - R 500,000 @ 25% cession
Income Level 3: R 500,001 - R 2,000,000 @ 40% cession

Member A: Salary R 150,000 → Level 1 → 10% ceded
Member B: Salary R 350,000 → Level 2 → 25% ceded
Member C: Salary R 800,000 → Level 3 → 40% ceded
```

## Support and Questions

For questions or issues related to the three-tier reinsurance structure:
- Review this documentation
- Check `backend/services/reinsurance_cession.go` for calculation logic
- Examine `backend/models/reinsurance_treaty.go` for data model
- Test with sample treaties in the UI before production use
