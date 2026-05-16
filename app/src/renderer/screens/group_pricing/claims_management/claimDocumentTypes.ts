export interface ClaimDocumentType {
  code: string
  name: string
  required: boolean
}

export const claimDocumentTypes: Record<string, ClaimDocumentType[]> = {
  GLA: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_deceased',
      name: 'Certified ID - Deceased',
      required: true
    },
    {
      code: 'certified_id_claimant',
      name: 'Certified ID - Claimant/Beneficiaries',
      required: true
    },
    {
      code: 'death_certificate',
      name: 'Death Certificate (BI-5)',
      required: true
    },
    {
      code: 'dha_notification',
      name: 'DHA-1663 Notification of Death',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'post_mortem',
      name: 'Post-mortem / Final BI-1680/1683',
      required: false
    }
  ],
  SGLA: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_deceased',
      name: 'Certified ID - Deceased',
      required: true
    },
    {
      code: 'certified_id_claimant',
      name: 'Certified ID - Claimant/Beneficiaries',
      required: true
    },
    {
      code: 'death_certificate',
      name: 'Death Certificate (BI-5)',
      required: true
    },
    {
      code: 'dha_notification',
      name: 'DHA-1663 Notification of Death',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: true
    },
    {
      code: 'relationship_proof',
      name: 'Proof of Relationship (Spouse/Child)',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'post_mortem',
      name: 'Post-mortem / Final BI-1680/1683',
      required: false
    }
  ],
  GFF: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_deceased',
      name: 'Certified ID - Deceased',
      required: true
    },
    {
      code: 'certified_id_claimant',
      name: 'Certified ID - Claimant/Beneficiaries',
      required: true
    },
    {
      code: 'death_certificate',
      name: 'Death Certificate (BI-5)',
      required: true
    },
    {
      code: 'dha_notification',
      name: 'DHA-1663 Notification of Death',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: true
    },
    {
      code: 'relationship_proof',
      name: 'Proof of Relationship (Spouse/Child)',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'post_mortem',
      name: 'Post-mortem / Final BI-1680/1683',
      required: false
    }
  ],
  PTD: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: false
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'attending_doctor_statement',
      name: "Attending Doctor's Statement (Disability/CI Report)",
      required: true
    },
    {
      code: 'specialist_report',
      name: 'Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)',
      required: false
    },
    {
      code: 'employer_duties_statement',
      name: 'Employer Statement of Duties / Job Description',
      required: true
    },
    {
      code: 'functional_capacity_assessment',
      name: 'Functional Capacity Assessment (FCE)',
      required: false
    },
    {
      code: 'occupational_therapist_report',
      name: 'Occupational Therapist Report',
      required: false
    }
  ],
  CI: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: false
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'attending_doctor_statement',
      name: "Attending Doctor's Statement (Disability/CI Report)",
      required: true
    },
    {
      code: 'specialist_report',
      name: 'Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)',
      required: true
    },
    {
      code: 'employer_duties_statement',
      name: 'Employer Statement of Duties / Job Description',
      required: false
    },
    {
      code: 'functional_capacity_assessment',
      name: 'Functional Capacity Assessment (FCE)',
      required: false
    },
    {
      code: 'occupational_therapist_report',
      name: 'Occupational Therapist Report',
      required: false
    }
  ],
  PHI: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: false
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'specialist_report',
      name: 'Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)',
      required: true
    },
    {
      code: 'employer_duties_statement',
      name: 'Employer Statement of Duties / Job Description',
      required: true
    },
    {
      code: 'functional_capacity_assessment',
      name: 'Functional Capacity Assessment (FCE)',
      required: false
    },
    {
      code: 'occupational_therapist_report',
      name: 'Occupational Therapist Report',
      required: false
    },
    {
      code: 'psychiatric_report',
      name: 'Psychiatric Report (if mental illness claim)',
      required: false
    },
    {
      code: 'income_loss_proof',
      name: 'Proof of Income Loss / Sick Leave Records',
      required: true
    }
  ],
  TTD: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    }
  ]
}
