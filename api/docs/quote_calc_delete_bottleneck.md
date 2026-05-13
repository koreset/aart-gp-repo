# Quote Calculation — Per-Category Delete/Insert Bottleneck

**Status:** open, deferred for future enhancement
**First diagnosed:** 2026-05-13
**Related earlier wins:** [api/services/group_pricing.go](../services/group_pricing.go) `GetPhiRate` / `GetReinsurancePhiRate` lookup_key index (migrations 20260513010000); `member_premium_schedules` composite index (migrations 20260513020000).

## Symptom

After the phi_rates and member_premium_schedules index changes landed, the dominant slow queries in a `CalculateGroupPricingQuote` run shifted to per-category `DELETE`s issued by [api/services/group_pricing.go:2148-2150](../services/group_pricing.go#L2148):

```
DELETE FROM `member_rating_results`    WHERE quote_id = '4' AND category = 'Administration'  -- 96 560 rows, ~16 s
DELETE FROM `member_premium_schedules` WHERE quote_id = '4' AND category = 'Administration'  -- 96 560 rows, ~7 s
DELETE FROM `member_rating_result_summaries` WHERE quote_id = ? AND category = ?            -- (smaller, not slow-flagged)
```

These run inside `calculateForCategory` after the in-memory rating math finishes, before the new result-set is bulk-inserted.

## Root cause (it is **not** the WHERE clause)

The WHERE filter on both DELETEs is already index-backed:

| Table                       | Index used                                           | Migration                                  |
|-----------------------------|------------------------------------------------------|--------------------------------------------|
| `member_rating_results`     | `idx_mrr_quote_category(quote_id, category)`         | `20260523000000_add_member_rating_indexes` |
| `member_premium_schedules`  | `idx_mps_quote_category_member(quote_id, category, member_name)` | `20260513020000_add_member_premium_schedules_index` |

Adding more indexes cannot speed this up — every secondary index makes the DELETE *slower*, because each deleted row's entry has to be removed from every index plus the InnoDB undo/redo log.

The cost is the inherent per-row DELETE machinery in InnoDB:

- `member_rating_results`: 3 secondary indexes → ~0.17 ms/row × 96 560 ≈ 16 s
- `member_premium_schedules`: 2 secondary indexes → ~0.07 ms/row × 96 560 ≈ 7 s

The pattern in `calculateForCategory` is "wipe this (quote_id, category) slice then re-insert it from scratch". On a 2-category quote, that's ~400 k row writes per recalculation.

## Side note: quoted `quote_id` in the slow-query log

The slow query log shows `WHERE quote_id = '4'` (quoted). `quoteId` is plumbed through `calculateForCategory(quoteId string, ...)` as a string. The `quote_id` column is BIGINT, so MySQL does an implicit string→int cast. Empirically the index still seeks (the `rows: 96560` count matches the actual row count rather than the whole table), so this is **not** the bottleneck — but it's a code smell worth tidying when the call signature is touched.

## Path A — UPDATE-in-place instead of DELETE+INSERT (recommended)

**Idea:** Hold the existing rows' primary keys through the calc, then `UPDATE … WHERE id = ?` only the columns whose values changed. Net effect:

- No 96k-row delete.
- No 96k-row bulk insert.
- Index maintenance only happens for rows whose *indexed* columns change (none of the ratings/premiums are indexed today, so most updates are cheap).

**Required work**

1. Add an `ID int gorm:"primary_key;autoIncrement"` to `MemberRatingResult` and `MemberPremiumSchedule` in [api/models/group_pricing.go](../models/group_pricing.go).
2. Migration per dialect: add `id BIGINT AUTO_INCREMENT PRIMARY KEY` (MySQL) / `BIGSERIAL` (PG) / `BIGINT IDENTITY` (MSSQL) and backfill. `member_premium_schedules` currently has no PK at all (see comment at [api/services/group_pricing.go:8449](../services/group_pricing.go#L8449)) — this also fixes that.
3. In `calculateForCategory`, after computing the new in-memory slice, build a `member_name|entry_date → existing_row_id` map from the current DB rows, then issue `UPDATE` statements. Insert rows that are new (member added since last calc); delete rows that disappeared (member exited). Most quote recalculations touch the same member set, so deletes/inserts shrink to near-zero.
4. Resolve the row-uniqueness footgun documented at [api/services/group_pricing.go:8449](../services/group_pricing.go#L8449). The current text warns that `(quote_id, category, member_name)` is not unique because movements produce multiple rows per member with different `entry_date`. With a real PK, this is moot — the ID disambiguates.

**Estimated payoff:** the 16 s + 7 s per category collapses to a few hundred ms; quote total drops by 60–80 s on a 2-category quote of this size.

**Risk:** medium. The delete+reinsert pattern is liberal — it absorbs schema drift, partial writes, and mid-run aborts. Switching to update-in-place needs an explicit reconciliation pass (which IDs disappeared? which appeared?) so an aborted run can be retried cleanly.

## Path B — investigate why "Administration" has 96 560 rows

Before any of Path A, it's worth one `GROUP BY` query to confirm whether 96 560 rows for *one category* is the intended row-count, or whether there's a duplication bug somewhere upstream that's blowing up the result set. The "Administration" category is a single scheme category — 96k rows implies either many members × many movement dates × many years, or a duplication regression.

```sql
SELECT
    COUNT(*)                                   AS rows,
    COUNT(DISTINCT member_name)                AS distinct_members,
    COUNT(DISTINCT year)                       AS distinct_years,
    COUNT(DISTINCT (entry_date))               AS distinct_entry_dates
FROM member_rating_results
WHERE quote_id = 4 AND category = 'Administration';
```

If `distinct_members × distinct_years` is much smaller than `rows`, there's a duplication bug and the right fix is to find that — Path A becomes unnecessary.

## Path C — partition by quote_id (do not recommend)

MySQL `PARTITION BY HASH(quote_id)` would make `TRUNCATE PARTITION` an instant delete, but MySQL partition keys must include every UNIQUE/PK column, partition pruning across joins is fragile, and the cross-dialect story (PG and MSSQL partitioning are very different) is painful. Not worth the operational cost for the gain.

## Decision

Park until the calc time becomes a customer-visible blocker. When picked up, run the diagnostic query from Path B first; if that doesn't reveal a duplication bug, do Path A. Don't add more indexes — they make this DELETE slower, not faster.
