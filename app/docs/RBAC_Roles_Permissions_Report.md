# AART Group Risk — Roles, Permissions & Governance Report

**Version:** 1.0 **Date:** 20 March 2026 **Prepared for:** AART v5.5.0 Group Pricing Module **Classification:** Internal — Confidential

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Role Definitions](#2-role-definitions)
   - 2.1 Internal Insurer/Administrator Roles
   - 2.2 External Portal Roles
   - 2.3 Specialist / Occasional Roles
3. [Module Access Matrix](#3-module-access-matrix)
4. [Claims Authority Levels & Escalation](#4-claims-authority-levels--escalation)
   - 4.1 Authority Tiers
   - 4.2 Automatic Escalation Triggers
   - 4.3 Turnaround Standards
5. [Segregation of Duties](#5-segregation-of-duties)
   - 5.1 Claims Processing
   - 5.2 Premium & Financial
   - 5.3 Underwriting & Pricing
   - 5.4 Bordereaux & Reinsurance
   - 5.5 System Administration
6. [Regulatory Requirements](#6-regulatory-requirements)
   - 6.1 POPIA
   - 6.2 Insurance Act 18 of 2017
   - 6.3 FAIS
   - 6.4 King IV
   - 6.5 TCF
   - 6.6 ASISA Guidelines
   - 6.7 Policyholder Protection Rules
7. [RBAC Design Principles for AART](#7-rbac-design-principles-for-aart)
8. [Appendix A — Permission Codes](#appendix-a--permission-codes)
9. [References](#references)

---

## 1. Executive Summary

This document defines the role-based access control (RBAC) framework for the AART Group Pricing application. It is informed by:

- Industry practices at SA group risk insurers (Liberty, Sanlam, Momentum, Old Mutual) and administrators (Alexander Forbes, NBC)
- Regulatory requirements from the FSCA, Insurance Act, POPIA, FAIS, and King IV
- The AART application's functional modules: Scheme Management, Quoting & Pricing, Member Management, Claims Management, Premium Management, Bordereaux Management, Reinsurance, and Reporting & Analytics

The framework defines **18 internal roles**, **3 external portal roles**, and **2 specialist roles**, supported by a 5-tier claims authority hierarchy, a comprehensive segregation of duties matrix, and regulatory alignment across 7 legislative instruments.

---

## 2. Role Definitions

### 2.1 Internal Insurer/Administrator Roles

#### System Administrator

- Manages user accounts, role assignments, system configuration, and audit logs
- Configures metadata tables (occupation classes, industry codes, benefit types)
- **No access** to transactional business data (claim amounts, member personal details) except via audit logs
- **Cannot** approve claims, process payments, or modify pricing

#### Scheme Administrator (Benefits Administrator)

- Day-to-day administration of group schemes: adding/removing members, salary updates, benefit changes, bulk uploads
- Processes bordereaux submissions from employers, validates member data against scheme rules
- Generates premium schedules from membership data
- **Cannot** approve claims, modify pricing parameters, or sign off financials

#### Senior Scheme Administrator / Team Leader

- All Scheme Administrator capabilities plus oversight of a team's work
- Can authorise member data changes above a threshold (e.g. salary changes >10%)
- Reviews and approves bordereaux submissions before sync to master register
- **Cannot** approve claims above junior threshold or modify rating tables

#### Underwriter / Pricing Actuary

- Creates and maintains quotes, sets rating factors, configures benefit structures (GLA, SGLA, PTD, CI, TTD, PHI, GFF)
- Runs experience analysis, adjusts loading factors, sets free cover limits
- Approves scheme-level underwriting decisions (new business, renewals, rate changes)
- **No access** to individual claims assessment; **no** payment processing capability
- Read-only access to claims analytics for pricing feedback

#### Actuarial Analyst

- Runs IBNR projections, experience analysis, IFRS 17 valuations (CSM, risk adjustment, PAA)
- Read-only access to claims data (aggregated and individual for reserving)
- Read-only access to premium and membership data
- **Cannot** modify schemes, members, or pricing directly; outputs feed into underwriter decisions
- Produces regulatory and management reports

#### Claims Administrator (Claims Registrar)

- Registers new claims: captures claimant details, event information, supporting documents
- Performs initial completeness checks (all required documentation received)
- **Cannot** assess, approve, decline, or pay claims — strictly intake and data capture
- Can update claim status to "registered" or "documents outstanding"

#### Claims Assessor (Junior)

- Assesses registered claims against policy terms, medical evidence, and scheme rules
- Recommends approval/decline with supporting rationale
- Authority to approve straightforward claims up to a defined limit (see Section 4)
- **Cannot** approve own registered claims (segregation of duties)
- **Cannot** process payments

#### Senior Claims Assessor

- All junior assessor capabilities with higher individual authority limit
- Reviews junior assessor recommendations for claims above junior threshold
- Handles complex claims (contested, multiple benefits, occupational disability)
- Can request additional medical opinions or forensic investigation
- **Cannot** approve claims above senior threshold; must escalate to Claims Manager

#### Claims Manager / Head of Claims

- Approves/declines claims up to management authority limit
- Oversees claims department operations, turnaround times, TCF compliance
- Refers claims above management limit to Claims Committee
- Signs off on claims analytics reports
- Can authorise ex-gratia payments within defined limits

#### Claims Committee Member

- Collective decision-making body for high-value or complex claims
- Approves/declines claims above Claims Manager authority
- Typical composition: Head of Claims, Chief Actuary, Chief Medical Officer, Legal Counsel, Chief Risk Officer
- Minutes must be recorded; decisions require quorum
- Refers to Board/Exco only for claims exceeding committee mandate

#### Finance Officer / Premium Accountant

- Manages premium invoicing, payment collection, reconciliation, arrears tracking
- Processes employer/broker statements
- Reconciles bank receipts against invoices
- **Cannot** modify scheme terms, member data, or claims decisions
- Can flag arrears for suspension action but **cannot** execute suspension

#### Finance Manager

- Approves payment runs (claims disbursements, reinsurance settlements)
- Signs off on premium reconciliation reports and arrears write-offs
- Reviews financial statements and technical accounts
- **Cannot** assess or approve claims; only authorises the financial execution of approved claims

#### Reinsurance Analyst

- Manages treaty administration: cession calculations, RI bordereaux generation
- Tracks reinsurer acceptances, recoveries, and technical accounts
- Prepares settlement statements
- Read-only access to claims data for recovery tracking
- **Cannot** modify claim decisions or scheme terms

#### Reinsurance Manager

- Approves outbound bordereaux submissions to reinsurers
- Negotiates treaty terms (outside system)
- Approves reinsurance settlement payments
- Reviews and signs off technical accounts

#### Compliance Officer

- Read-only access to all modules for monitoring purposes
- Reviews audit trails, user access logs, data access patterns
- Monitors TCF outcome compliance, POPIA adherence, FAIS requirements
- Can flag transactions for review but **cannot** modify operational data
- Produces compliance reports for FSCA and Prudential Authority

#### Risk Officer

- Read-only access to all modules
- Reviews risk registers, operational risk events, control effectiveness
- Monitors claims authority limit adherence and segregation of duties compliance
- Reports to Board Risk Committee

#### Internal Auditor

- Read-only access to all modules and full audit trail
- Conducts periodic reviews of controls, segregation of duties, authority limits
- **Cannot** modify any operational data
- Reports to Audit Committee (independent of management)

#### IT Administrator / System Support

- Technical system maintenance, database administration, backup management
- Access to system configuration but **not** to business data in production (or only via controlled, logged access)
- **Cannot** process transactions, approve claims, or view member personal information in bulk

### 2.2 External Portal Roles

#### Broker / Financial Adviser

- Views schemes under their agency code only
- Can submit new business quotes, view quote results, request renewals
- Views premium statements and commission reports for their book
- Views claims status (registered, in assessment, approved, declined) — **no** claim amounts or medical details
- Can submit claims on behalf of employer clients
- **Cannot** modify scheme terms, assess claims, or access other brokers' data
- Must be FAIS-compliant (Fit & Proper, RE exams, CPD)

#### Employer HR Administrator

- Views and manages members for their scheme(s) only
- Submits bordereaux (joiner/leaver/amendment files)
- Views premium schedules and invoices for their scheme
- Can download statements
- Views claim status for their scheme members (limited detail)
- **Cannot** access pricing, other employers' data, or claims assessment details
- **Cannot** modify benefit structures or contribution rates

#### Employer Finance Contact

- Views invoices and payment history for their scheme only
- Can record proof of payment
- Views arrears status
- **No access** to individual member data, claims, or pricing

### 2.3 Specialist / Occasional Roles

#### Medical Adviser / Chief Medical Officer

- Consulted on complex or contested disability/critical illness claims
- Access to individual claim medical documentation (POPIA: purpose-limited)
- Provides medical opinions; **does not** approve/decline claims
- Claims Committee member

#### Forensic Investigator

- Access to specific flagged claims only (not broad access)
- Views claim documentation, member history, and related scheme data
- Produces investigation reports
- **Cannot** approve/decline claims

---

## 3. Module Access Matrix

Legend: **CRUD** = Create/Read/Update/Delete | **R** = Read-only | **—** = No access | **(own)** = Own book/scheme only

| Module | Sys Admin | Scheme Admin | Sr Scheme Admin | Underwriter | Actuary | Claims Registrar | Jr Assessor | Sr Assessor | Claims Mgr | Finance Officer | Finance Mgr | RI Analyst | RI Mgr | Compliance | Broker | Employer HR | Employer Fin |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| **Scheme Management** | Config | CRUD | CRUD+Approve | R+Quote | R | — | R | R | R | R | R | R | R | R | R (own) | R (own) | — |
| **Quoting & Pricing** | Config | — | — | CRUD | R+Run | — | — | — | — | — | — | — | — | R | Submit+R | — | — |
| **Member Management** | Config | CRUD | CRUD+Approve | R | R | R | R | R | R | — | — | — | — | R | R (own) | CRUD (own) | — |
| **Claims — Register** | — | — | — | — | — | CRUD | R | R | R | — | — | — | — | R | Submit | R (own) | — |
| **Claims — Assess** | — | — | — | — | R | — | CRUD | CRUD | CRUD | — | — | — | — | R | — | — | — |
| **Claims — Approve** | — | — | — | — | — | — | ≤Tier 1 | ≤Tier 2 | ≤Tier 3 | — | — | — | — | R | — | — | — |
| **Claims — Committee** | — | — | — | — | — | — | — | — | Refer | — | — | — | — | R | — | — | — |
| **Claims — Pay (Initiate)** | — | — | — | — | — | — | — | — | Initiate | — | — | — | — | R | — | — | — |
| **Claims — Pay (Authorise)** | — | — | — | — | — | — | — | — | — | — | Approve | — | — | R | — | — | — |
| **Claims — ACB Generate** | — | — | — | — | — | — | — | — | Generate | R | Approve | — | — | R | — | — | — |
| **Claims — Reconcile** | — | — | — | — | — | — | — | — | — | CRUD | R | — | — | R | — | — | — |
| **Claims Analytics** | — | — | — | R | CRUD | — | R | R | CRUD | R | R | R | R | R | — | — | — |
| **Premium Schedules** | Config | Generate | Approve | R | R | — | — | — | — | CRUD | R | — | — | R | R (own) | R (own) | R (own) |
| **Invoicing** | Config | R | R | — | — | — | — | — | — | CRUD | Approve | — | — | R | R (own) | R (own) | R (own) |
| **Payments & Recon** | — | — | — | — | — | — | — | — | — | CRUD | Approve | — | — | R | — | — | R (own) |
| **Arrears** | — | Action | Approve | — | — | — | — | — | — | CRUD | Approve | — | — | R | R (own) | R (own) | — |
| **Statements** | — | — | — | — | — | — | — | — | — | Generate | Approve | — | — | R | R (own) | — | R (own) |
| **Bordereaux Inbound** | Config | CRUD | Approve | — | R | — | — | — | — | — | — | — | — | R | — | Submit (own) | — |
| **Bordereaux Outbound** | Config | — | — | — | R | — | — | — | — | — | — | Generate | Approve+Submit | R | — | — | — |
| **RI Tracking** | — | — | — | — | R | — | — | — | — | R | R | CRUD | Approve | R | — | — | — |
| **Claim Notifications** | — | — | — | — | — | — | — | — | R | — | — | CRUD | Approve | R | — | — | — |
| **RI Treaties** | Config | — | — | R | R | — | — | — | — | R | R | CRUD | Approve | R | — | — | — |
| **RI Settlement** | — | — | — | — | R | — | — | — | — | R | Approve | CRUD | Approve | R | — | — | — |
| **Experience Analysis** | — | — | — | Run+R | Run+R | — | — | — | R | — | — | R | R | R | — | — | — |
| **IBNR Projections** | — | — | — | R | Run+R | — | — | — | R | R | R | R | R | R | — | — | — |
| **IFRS 17 Reporting** | — | — | — | R | Run+R | — | — | — | — | R | R | R | R | R | — | — | — |
| **User Management** | CRUD | — | — | — | — | — | — | — | — | — | — | — | — | R | — | — | — |
| **Audit Logs** | R | — | — | — | — | — | — | — | — | — | — | — | — | R | — | — | — |
| **System Config** | CRUD | — | — | — | — | — | — | — | — | — | — | — | — | R | — | — | — |

---

## 4. Claims Authority Levels & Escalation

### 4.1 Authority Tiers

All amounts are per-claim, per-life. These are representative thresholds — each insurer sets their own via Board-approved Delegation of Authority (DoA).

| Tier | Role | GLA / SGLA | Disability (PTD/CI) | Income Benefits (TTD/PHI) | Funeral (GFF) |
| --- | --- | --- | --- | --- | --- |
| 1 | Claims Assessor (Junior) | Up to R500,000 | Up to R250,000 lump sum | Up to R50,000/month | Up to R75,000 |
| 2 | Senior Claims Assessor | R500,001 – R2,000,000 | R250,001 – R1,000,000 | R50,001 – R150,000/month | R75,001 – R150,000 |
| 3 | Claims Manager | R2,000,001 – R10,000,000 | R1,000,001 – R5,000,000 | R150,001 – R500,000/month | All funeral claims |
| 4 | Claims Committee | Above R10,000,000 | Above R5,000,000 | Above R500,000/month | N/A |
| 5 | Board / Exco | Exceeds free cover limit or reinsurance retention | Exceptional circumstances | Exceptional circumstances | N/A |

### 4.2 Automatic Escalation Triggers

The following conditions require escalation to the next authority tier **regardless of claim amount**:

1. **Early claims** — event within 12 months of cover inception or reinstatement
2. **Key person / principal member** — claimant has scheme influence (e.g. trustee, employer director)
3. **Related party** — claimant is related to an employee of the insurer or administrator
4. **Suspected fraud or material non-disclosure** — flagged during assessment or investigation
5. **Suicide within exclusion period** — typically 24 months from inception
6. **Contested or disputed claim** — employer or beneficiary disputes the outcome
7. **Ex-gratia consideration** — claim falls outside strict policy terms but compassionate grounds exist
8. **Reinsurer notification required** — claim exceeds automatic retention under treaty terms
9. **Multiple simultaneous claims** — e.g. catastrophe event affecting multiple members of one scheme

### 4.3 Turnaround Standards

Aligned with TCF Outcome 6 and Policyholder Protection Rules:

| Claim Type | Decision Target | Payment Target |
| --- | --- | --- |
| Straightforward (all docs received) | 5 business days | 3 business days after approval |
| Complex (medical opinion required) | 15 business days | 3 business days after approval |
| Claims Committee referral | Next scheduled meeting (fortnightly/monthly) | 3 business days after committee decision |
| Reinsurer-notifiable | Subject to reinsurer response time | 3 business days after reinsurer acknowledgement |

---

## 5. Segregation of Duties

### 5.1 Claims Processing

| Function A | Function B | Rationale |
| --- | --- | --- |
| Register claim | Assess or approve the same claim | Prevents registration of fictitious claims |
| Assess claim | Approve the same claim (above assessor's authority) | Four-eyes principle on material financial decisions |
| Approve claim | Authorise payment for the same claim | Prevents payment of unapproved or manipulated claims |
| Authorise payment | Execute bank payment (generate/upload ACB file) | Prevents unauthorized fund transfers |
| Assess claim | Modify member or beneficiary data | Prevents manipulation of benefit calculation inputs |

### 5.2 Premium & Financial

| Function A | Function B | Rationale |
| --- | --- | --- |
| Generate premium schedule | Approve or finalize the schedule | Prevents incorrect billing going undetected |
| Generate invoice | Record payment against the same invoice | Prevents fictitious payment recording |
| Record payment | Reconcile the same payment | Prevents concealment of misappropriated funds |
| Initiate arrears suspension | Approve the suspension | Prevents arbitrary or malicious member exclusion |
| Generate employer statement | Approve write-off on the same statement | Prevents unauthorized premium write-offs |

### 5.3 Underwriting & Pricing

| Function A | Function B | Rationale |
| --- | --- | --- |
| Set rating tables or loading factors | Approve a quote for a specific client using those tables | Prevents self-serving pricing manipulation |
| Generate a quote | Accept or bind a scheme on that quote | Prevents binding without independent commercial review |
| Modify free cover limits | Approve individual underwriting decisions | Prevents circumventing medical underwriting requirements |

### 5.4 Bordereaux & Reinsurance

| Function A | Function B | Rationale |
| --- | --- | --- |
| Generate outbound bordereaux | Approve or submit the bordereaux to the reinsurer | Prevents submission of incorrect or manipulated data |
| Review bordereaux data | Final approval of the same bordereaux | Prevents single-person sign-off on material data |
| Calculate cession amounts | Approve reinsurance settlement payments | Prevents incorrect reinsurance payments |
| Record reinsurer recovery | Reconcile the RI technical account | Prevents concealment of recovery variances |

### 5.5 System Administration

| Function A | Function B | Rationale |
| --- | --- | --- |
| Create or modify user accounts and roles | Perform any business transactions | Prevents self-granting of elevated access to commit fraud |
| Modify system configuration or metadata | Approve business transactions affected by those settings | Prevents system manipulation for personal gain |
| Access audit logs (read) | Modify or delete audit log data | Audit trail integrity — logs must be immutable |

---

## 6. Regulatory Requirements

### 6.1 POPIA (Protection of Personal Information Act, 2013)

| Requirement | System Impact |
| --- | --- |
| **Purpose limitation** | Users may only access personal information necessary for their specific function. Finance officers do not need medical claim documentation. Claims assessors do not need premium payment details. |
| **Minimality** | Field-level access control where feasible — masking ID numbers, hiding medical details from non-claims staff, restricting salary visibility. |
| **Security safeguards** | Role-based access control is explicitly expected. Access logs must be maintained and auditable. |
| **Data subject access** | Members/employers can request what data is held about them. System must support data subject access requests (DSARs). |
| **Cross-border transfers** | Reinsurer data sent offshore (e.g. to RGA, Munich Re, Swiss Re) must comply with Section 72 conditions. Reinsurance module users must be aware. |
| **Information Officer** | Organisation must designate an Information Officer and Deputy Information Officers. System should record who accesses what personal data. |
| **Retention limitation** | Claims and member data must have defined retention periods. System must support archiving and deletion policies. |

### 6.2 Insurance Act 18 of 2017

| Requirement | System Impact |
| --- | --- |
| **Governance framework** | Clear organisational structure, well-defined lines of responsibility, effective risk management, and internal controls — all reflected in RBAC design. |
| **Board oversight** | Board is ultimately responsible for governance, including delegation of authority structures that define claims authority tiers. |
| **Outsourcing / binder agreements** | If the administrator operates under a binder agreement, the insurer must maintain oversight capability — requiring read-only access or regular data extracts from the administrator's system. |
| **Claims management framework** | PPR requires documented procedures, timeframes, escalation, and complaint handling — all enforced via system workflow and audit trail. |

### 6.3 FAIS (Financial Advisory and Intermediary Services Act, 2002)

| Requirement | System Impact |
| --- | --- |
| **Fit & Proper requirements** | Key Individuals and Representatives must meet competency requirements. System should track user qualifications and certifications. |
| **Broker licensing** | Only FAIS-licensed brokers may access broker portal. System should validate FSP numbers at registration. |
| **Advice records** | If the system captures advice (e.g. quote recommendations), it must be retained for 5 years minimum. |
| **Conflict of interest** | Broker commission visibility must be restricted — brokers should not see competitor commission structures. |

### 6.4 King IV Code on Corporate Governance

| Principle | System Impact |
| --- | --- |
| **Three lines of defence** | First line (operations) has transactional access; second line (risk/compliance) has monitoring/read access plus flagging; third line (internal audit) has full read access and audit trail. Each line must have appropriate but **distinct** system access. |
| **Board committees** | Audit Committee, Risk Committee, and Social & Ethics Committee should receive system-generated reports but do not require direct system access. |
| **Delegation of authority** | Must be documented, Board-approved, and enforced by the system's claims authority hierarchy. |

### 6.5 TCF (Treating Customers Fairly)

| Outcome | System Impact |
| --- | --- |
| **Outcome 5** (Products perform as expected) | Scheme administrators and underwriters must have access to ensure benefits are correctly configured and match policy terms. |
| **Outcome 6** (No unreasonable post-sale barriers) | Claims registration must be accessible. Employer portals must allow timely submission. System must track and report claims turnaround times. |
| **Audit trail** | All claim status changes must be logged with timestamps, user IDs, and rationale to demonstrate fair treatment. |

### 6.6 ASISA Guidelines

| Guideline | System Impact |
| --- | --- |
| **Claims experience data** | Insurers must provide 5–7 years of claims experience data at scheme level for broker/employer review at renewal. System must support this reporting. |
| **Group Risk Committee standards** | ASISA publishes guidelines on claims experience sharing and group risk market conduct that may affect what data brokers can access. |

### 6.7 Long-term Insurance PPR (Policyholder Protection Rules, 2017)

| Requirement | System Impact |
| --- | --- |
| **Claims management framework** | Reasonable time for policyholders to institute claims. Prominent disclosure of time limitations. Procedures for keeping claimants informed. Internal escalation and complaint mechanisms. |
| **System implications** | Claim status must be visible to authorised parties (employer HR, broker) without exposing assessment details. Automated notifications should be supported. |

---

## 7. RBAC Design Principles for AART

Based on the analysis above, the following principles should guide AART's RBAC implementation:

1. **Minimum 18–20 distinct roles** covering internal operations, external portals, oversight functions, and specialist access.

2. **Hierarchical claims authority** with at least 4 tiers (Junior Assessor → Senior Assessor → Claims Manager → Claims Committee) plus Board escalation for exceptional cases.

3. **Strict SoD enforcement** — the system must prevent the same user from performing both sides of any segregation-of-duties pair listed in Section 5. This should be enforced at the application level, not just by policy.

4. **Data isolation for external users** — brokers see only their book (filtered by agency/FSP code); employers see only their scheme (filtered by scheme ID). No cross-contamination.

5. **Purpose-limited medical data access** — only Claims Assessors, Claims Managers, Medical Advisers, and Forensic Investigators may view claim medical documentation. All other roles see claim status only.

6. **Immutable audit trail** — no role, including System Administrator, can modify or delete audit logs. Compliance and Internal Audit have read access. Logs must include user ID, timestamp, action, affected record, and before/after values.

7. **POPIA field-level masking** — ID numbers, medical information, and salary data must be masked for roles that do not require them. Masking should be enforced at the API layer, not just the UI.

8. **Entitlement-driven UI** — aligns with AART's existing pattern where the Pinia store holds entitlements and route guards plus component `v-if` conditions check entitlements before rendering. Each permission should map to a discrete entitlement code.

9. **Claims authority limits configurable** — authority thresholds should be stored as system configuration (not hard-coded) so each client can set their own DoA-aligned limits.

10. **Dual authorisation for financial actions** — payment runs (ACB file generation), arrears write-offs, and reinsurance settlements require initiation by one role and approval by a different, more senior role.

---

## Appendix A — Permission Codes

Suggested permission code structure for AART, using the format `module:action`:

### Scheme Management

| Code                      | Description                                 |
| ------------------------- | ------------------------------------------- |
| `schemes:view`            | View scheme list and details                |
| `schemes:create`          | Create new schemes                          |
| `schemes:edit`            | Edit scheme details                         |
| `schemes:delete`          | Delete/archive schemes                      |
| `schemes:approve_changes` | Approve member data changes above threshold |

### Quoting & Pricing

| Code                    | Description                              |
| ----------------------- | ---------------------------------------- |
| `quotes:view`           | View quotes                              |
| `quotes:create`         | Create new quotes                        |
| `quotes:edit_rating`    | Modify rating tables and loading factors |
| `quotes:approve`        | Approve/bind a quote                     |
| `quotes:run_experience` | Run experience analysis                  |

### Member Management

| Code                  | Description                                 |
| --------------------- | ------------------------------------------- |
| `members:view`        | View member list and details                |
| `members:create`      | Add new members                             |
| `members:edit`        | Edit member details                         |
| `members:bulk_upload` | Perform bulk member uploads                 |
| `members:view_salary` | View salary information (masked by default) |

### Claims Management

| Code                           | Description                       |
| ------------------------------ | --------------------------------- |
| `claims:view`                  | View claim list and status        |
| `claims:register`              | Register new claims               |
| `claims:assess`                | Create claim assessments          |
| `claims:approve_tier1`         | Approve claims up to Tier 1 limit |
| `claims:approve_tier2`         | Approve claims up to Tier 2 limit |
| `claims:approve_tier3`         | Approve claims up to Tier 3 limit |
| `claims:refer_committee`       | Refer claims to Claims Committee  |
| `claims:committee_decide`      | Record Claims Committee decisions |
| `claims:decline`               | Decline claims (within authority) |
| `claims:view_medical`          | View medical documentation        |
| `claims:view_analytics`        | View claims analytics dashboard   |
| `claims:request_info`          | Request additional information    |
| `claims:request_investigation` | Request forensic investigation    |

### Claims Payments

| Code                              | Description                |
| --------------------------------- | -------------------------- |
| `claims_pay:create_schedule`      | Create payment schedules   |
| `claims_pay:generate_acb`         | Generate ACB files         |
| `claims_pay:authorise_payment`    | Authorise payment runs     |
| `claims_pay:upload_response`      | Upload bank response files |
| `claims_pay:reconcile`            | Reconcile bank responses   |
| `claims_pay:retry_failed`         | Retry failed payments      |
| `claims_pay:manage_bank_profiles` | Create/edit bank profiles  |

### Premium Management

| Code                          | Description                    |
| ----------------------------- | ------------------------------ |
| `premiums:view`               | View premium data              |
| `premiums:generate_schedule`  | Generate premium schedules     |
| `premiums:approve_schedule`   | Approve/finalize schedules     |
| `premiums:generate_invoice`   | Generate invoices              |
| `premiums:record_payment`     | Record premium payments        |
| `premiums:reconcile`          | Perform premium reconciliation |
| `premiums:manage_arrears`     | Manage arrears actions         |
| `premiums:approve_suspension` | Approve arrears suspensions    |
| `premiums:generate_statement` | Generate statements            |
| `premiums:approve_writeoff`   | Approve premium write-offs     |

### Bordereaux

| Code                           | Description                                |
| ------------------------------ | ------------------------------------------ |
| `bordereaux:view`              | View bordereaux data                       |
| `bordereaux:submit_inbound`    | Submit inbound employer bordereaux         |
| `bordereaux:process_inbound`   | Process and sync inbound submissions       |
| `bordereaux:approve_inbound`   | Approve inbound submissions                |
| `bordereaux:generate_outbound` | Generate outbound bordereaux               |
| `bordereaux:approve_outbound`  | Approve outbound bordereaux for submission |
| `bordereaux:submit_outbound`   | Submit bordereaux to reinsurer             |

### Reinsurance

| Code | Description |
| --- | --- |
| `reinsurance:view` | View reinsurance data |
| `reinsurance:manage_treaties` | Create/edit treaty configurations |
| `reinsurance:approve_treaties` | Approve treaty terms |
| `reinsurance:manage_cessions` | Calculate and manage cessions |
| `reinsurance:manage_recoveries` | Track reinsurer recoveries |
| `reinsurance:manage_settlement` | Prepare settlement statements |
| `reinsurance:approve_settlement` | Approve settlement payments |
| `reinsurance:manage_notifications` | Manage claim notifications to reinsurers |

### Reporting & Analytics

| Code                     | Description                 |
| ------------------------ | --------------------------- |
| `reports:view`           | View standard reports       |
| `reports:run_ibnr`       | Run IBNR projections        |
| `reports:run_ifrs17`     | Run IFRS 17 valuations      |
| `reports:run_experience` | Run experience analysis     |
| `reports:export`         | Export reports to Excel/PDF |

### Administration

| Code | Description |
| --- | --- |
| `admin:manage_users` | Create/edit/delete user accounts |
| `admin:manage_roles` | Create/edit role definitions and permissions |
| `admin:view_audit_logs` | View audit trail |
| `admin:system_config` | Modify system configuration |
| `admin:manage_authority_limits` | Configure claims authority thresholds |

---

## References

1. FSCA — Financial Sector Conduct Authority. https://www.fsca.co.za/
2. Insurance Act 18 of 2017. https://www.gov.za/documents/acts/insurance-act-18-2017
3. Long-Term Insurance Act: Policyholder Protection Rules (2017). Government Gazette No. 41321.
4. Protection of Personal Information Act (POPIA), 2013. https://popia.co.za/
5. King IV Report on Corporate Governance for South Africa (2016). Institute of Directors in Southern Africa.
6. FAIS Act — Financial Advisory and Intermediary Services Act, 2002. FSCA Fit & Proper Requirements.
7. ASISA — Association for Savings and Investment South Africa. Standards, Guidelines and Codes. https://asisa.org.za/
8. ASISA Updated Guidelines on Claims Experience for Group Schemes.
9. TCF — Treating Customers Fairly. FSCA Regulatory Framework.
10. IMF South Africa Financial Sector Assessment — Insurance Regulation and Supervision.
11. Board Notice 158 of 2014 — Governance and Risk Management Framework for Insurers (FSCA).

---

_This document should be reviewed and updated as regulatory requirements evolve and as AART's feature set expands. Claims authority limits in Section 4 are representative — each client deployment should configure their own Board-approved thresholds._
