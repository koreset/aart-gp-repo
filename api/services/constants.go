package services

const DUMMY_CODE = "COMP_01"
const TASK_PENDING = "pending"
const TASK_ACTIVE = "active"
const TASK_COMPLETE = "complete"
const MpLimitValue = 1

const MainMember = "MM"
const Spouse = "SP"
const Child = "CH"
const Parent = "PAR"

const Male = "M"
const Female = "F"

const Retrenchment = "Retrenchment"
const Health = "Health"
const Death = "Death"
const PermanentDisability = "Permanent Disability"
const AccidentalDeath = "Accidental Death"
const Lapse = "Lapse"
const TemporaryDisability = "Temporary Disability"
const CriticalIllness = "Critical Illness"
const Sickness = "Sickness"
const Maturity = "Maturity"
const FullyPaidUp = "Fully Paid Up"
const SemiPaidUp = "Semi Paid Up"

// Yield Curve Basis
const Current = "Current"
const LockedInRates = "Locked in Rates"

// Risk Drivers
const DISCOUNTED_DEATH_OUTGO = "DISCOUNTED_DEATH_OUTGO"
const DISCOUNTED_MORBIDITY_OUTGO = "DISCOUNTED_MORBIDITY_OUTGO"
const DISCOUNTED_EXPENSE_OUTGO = "DISCOUNTED_EXPENSE_OUTGO"
const DISCOUNTED_ANNUITY_INCOME = "DISCOUNTED_ANNUITY_INCOME"
const DISCOUNTED_OUTGO = "DISCOUNTED_OUTGO"
const DISCOUNTED_SURRENDER_OUTGO = "DISCOUNTED_SURRENDER_OUTGO"
const SUM_AT_RISK = "SUM_AT_RISK"
const RESERVES = "RESERVES"

// CommissionType
const Initial = "1"
const Renewal = "2"
const Hybrid = "3"

const defaultdecrementPrecision = 25
const defaultPrecision = 7
const CsmDefaultPrecision = 5
const AccountingPrecision = 2

// Retrenchment
const RetrMaxYearDimension = 6 // maximum year in view value. Set to retr_max_dimension if greater

// Lapse Table
const LapseMaxMonthDimension = 60 // maximum month in view. Set to the max value if greater

// Aggregation Period
const ScopedAggregatedProjectionPeriod = 72 // maximum projection month
const MaxProjectionMonthSap = 61            // maximum projection month of results read in CSM engine.should be less than ScopedAggregatedProjectionPeriod
const MaxCsmProjectionPeriodYears = 5       // maximum projection month of results read in CSM engine. should be less than MaxProjectionMonthSap/12
const LICAggregatedProjectionPeriod = 24

// BalanaceSheetRecord TransitionType
const PostTransition = "PostTransition"
const FullyRetrospective = "FullyRetrospective"
const ModifiedRetrospective = "ModifiedRetrospective"
const FairValue = "FairValue"

// gmm model point status
const IF = "IF"
const NB = "NB"

// finance
const Exit = "Exit"

// Portfolio Premium Earning Pattern
const PassageofTime = "passageoftime"
const SpecifiedbyUser = "specifiedbyuser"
const DailyPassageofTime = "dailypassageoftime"

// Coverage units discounting
const UndiscountedCoverageUnits = "UndiscountedCoverageUnits"
const DiscountedCoverageUnits = "DiscountedCoverageUnits"

// ibnr constant variables
const AccidentYear = "AccidentYear"
const AccidentMonth = "AccidentMonth"
const EarnedPremium = "EarnedPremium"
const ChainLadder = "chain-ladder"
const BornHuetterFerguson = "bornhuetter-ferguson"
const ChainLadderAverageCostperClaim = "chain-ladder-average-cost-per-claim"
const CapeCod = "cape-cod"
const ChainLadderBF = "chain-ladder-bf"
const ChainLadderCapeCod = "chain-ladder-cape-cod"
const ResidualNumber = "ResidualNumber"
const Residual = "Residual"
const RandomResidual = "RandomResidual"
const Monthly = "monthly"
const Annual = "annual"
const Quarterly = "quarterly"

//ibnr column reference numbers

const TriangleVariableCount = 8
const FactorVariableCount = 5
const VariableCountDiff = 3
const ibnrRatioPrecision = 7

// portfolios
const Undiscounted = "undiscounted"
const Discounted = "discounted"

// reinsurance
const Direct = "direct"
const ProportionalReinsurance = "proportionalreinsurance"
const NonProportionalReinsurance = "nonproportionalreinsurance"
const Inward = "inward"
const Outward = "outward"

// reinsuranceparameters
const ReinsuranceInwardOutward = "ReinsuranceInwardOutward"
const UlrLowerboundRate = "UlrLowerboundRate"
const UlrUpperboundRate = "UlrUpperboundRate"
const UlrLowerboundCommissionRate = "UlrLowerboundCommissionRate"
const UlrUpperboundCommissionRate = "UlrUpperboundCommissionRate"
const SlidingScaleMinRate = "SlidingScaleMinRate"
const SlidingScaleMaxRate = "SlidingScaleMaxRate"
const ProvisionalCommissionRate = "ProvisionalCommissionRate"
const ProfitCommissionModel = "ProfitCommissionModel"
const ProfitCommissionRate = "ProfitCommissionRate"
const Cedant = "cedant"
const Reinsurer = "reinsurer"

// IBNR Distribution

const Normal = "normal"
const Lognormal = "lognormal"
const Gamma = "gamma"
const Resampling = "resampling"

const IFRS17 = "IFRS17"

// Group Pricing Schemes
const IN_FORCE = "In Force"
const LAPSED = "Lapsed"
const EXPIRED = "Expired"
const TERMINATED = "Terminated"
const NTU = "Not Taken Up"
const CANCELLED = "Cancelled"
const REINSTATED = "Reinstated"
const SUSPENDED = "Suspended"
const PENDING_ACTIVATION = "Pending Activation"
const QUOTED = "Quoted"
const DECLINED = "Declined"
const ACCEPTED = "Accepted"

//colaims

const Pending = "pending"
const Notified = "notified"
