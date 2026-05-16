import Api from '@/renderer/api/Api'

export type FraudRiskLevel = 'low' | 'medium' | 'high' | 'critical'

export interface FraudFeatureSpec {
  name: string
  kind: 'numeric' | 'string' | 'bool'
  description: string
  used_by_glm: boolean
  used_by_rules: boolean
  choices?: string[]
}

export interface FraudRiskModel {
  id: number
  intercept: number
  coefficients: Record<string, number> | null
  trained_at: string | null
  trained_by: string
  sample_size: number
  positive_count: number
  auc: number
  updated_at: string
}

export type FraudRuleOp = 'eq' | 'ne' | 'gt' | 'gte' | 'lt' | 'lte' | 'in'

export interface FraudRuleLeaf {
  field: string
  op: FraudRuleOp
  value: number | string | boolean | Array<number | string>
}

/* eslint-disable no-use-before-define */
// The recursive group/node pair below is mutually referential by design;
// TypeScript hoists interfaces and types so this is safe at runtime.
export interface FraudRuleGroup {
  all?: FraudRuleNode[]
  any?: FraudRuleNode[]
}

export type FraudRuleNode = FraudRuleLeaf | FraudRuleGroup
/* eslint-enable no-use-before-define */

export interface FraudRiskRule {
  id: number
  name: string
  description: string
  conditions: FraudRuleNode
  risk_level: FraudRiskLevel
  priority: number
  enabled: boolean
  updated_by: string
  updated_at: string
  created_at: string
}

export interface MatchedRule {
  id: number
  name: string
  risk_level: FraudRiskLevel
}

export interface FraudCheckResult {
  glm_score: number
  glm_band: FraudRiskLevel
  matched_rule: MatchedRule | null
  final_risk_level: FraudRiskLevel
  rationale: string
  features: Record<string, number | string>
  assessment_id: number
}

export interface RefitResult {
  sample_size: number
  positive_count: number
  auc: number
  intercept: number
  coefficients: Record<string, number>
}

export default {
  runFraudCheck(claimId: number) {
    return Api.post<FraudCheckResult>(
      `/group-pricing/claims/${claimId}/fraud-check`
    )
  },
  getModel() {
    return Api.get<FraudRiskModel>('/group-pricing/fraud-risk/model')
  },
  refitModel() {
    return Api.post<RefitResult>('/group-pricing/fraud-risk/refit')
  },
  getFeatureCatalogue() {
    return Api.get<FraudFeatureSpec[]>('/group-pricing/fraud-risk/features')
  },
  listRules() {
    return Api.get<FraudRiskRule[]>('/group-pricing/fraud-risk/rules')
  },
  createRule(rule: Partial<FraudRiskRule>) {
    return Api.post<FraudRiskRule>('/group-pricing/fraud-risk/rules', rule)
  },
  updateRule(ruleId: number, rule: Partial<FraudRiskRule>) {
    return Api.put<FraudRiskRule>(
      `/group-pricing/fraud-risk/rules/${ruleId}`,
      rule
    )
  },
  deleteRule(ruleId: number) {
    return Api.delete(`/group-pricing/fraud-risk/rules/${ruleId}`)
  }
}
