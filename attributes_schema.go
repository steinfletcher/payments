package acme

const AttributesSchema = `{
  "definitions": {},
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "http://example.com/root.json",
  "type": "object",
  "title": "The Root Schema",
  "required": [
    "amount",
    "beneficiary_party",
    "charges_information",
    "currency",
    "debtor_party",
    "end_to_end_reference",
    "fx",
    "numeric_reference",
    "payment_id",
    "payment_purpose",
    "payment_scheme",
    "payment_type",
    "processing_date",
    "reference",
    "scheme_payment_sub_type",
    "scheme_payment_type",
    "sponsor_party"
  ],
  "properties": {
    "amount": {
      "$id": "#/properties/amount",
      "type": "string",
      "title": "The Amount Schema",
      "default": "",
      "examples": [
        "100.21"
      ],
      "pattern": "^(.*)$"
    },
    "beneficiary_party": {
      "$id": "#/properties/beneficiary_party",
      "type": "object",
      "title": "The Beneficiary_party Schema",
      "required": [
        "account_name",
        "account_number",
        "account_number_code",
        "account_type",
        "address",
        "bank_id",
        "bank_id_code",
        "name"
      ],
      "properties": {
        "account_name": {
          "$id": "#/properties/beneficiary_party/properties/account_name",
          "type": "string",
          "title": "The Account_name Schema",
          "default": "",
          "examples": [
            "W Owens"
          ],
          "pattern": "^(.*)$"
        },
        "account_number": {
          "$id": "#/properties/beneficiary_party/properties/account_number",
          "type": "string",
          "title": "The Account_number Schema",
          "default": "",
          "examples": [
            "31926819"
          ],
          "pattern": "^(.*)$"
        },
        "account_number_code": {
          "$id": "#/properties/beneficiary_party/properties/account_number_code",
          "type": "string",
          "title": "The Account_number_code Schema",
          "default": "",
          "examples": [
            "BBAN"
          ],
          "pattern": "^(.*)$"
        },
        "account_type": {
          "$id": "#/properties/beneficiary_party/properties/account_type",
          "type": "integer",
          "title": "The Account_type Schema",
          "default": 0,
          "examples": [
            0
          ]
        },
        "address": {
          "$id": "#/properties/beneficiary_party/properties/address",
          "type": "string",
          "title": "The Address Schema",
          "default": "",
          "examples": [
            "1 The Beneficiary Localtown SE2"
          ],
          "pattern": "^(.*)$"
        },
        "bank_id": {
          "$id": "#/properties/beneficiary_party/properties/bank_id",
          "type": "string",
          "title": "The Bank_id Schema",
          "default": "",
          "examples": [
            "403000"
          ],
          "pattern": "^(.*)$"
        },
        "bank_id_code": {
          "$id": "#/properties/beneficiary_party/properties/bank_id_code",
          "type": "string",
          "title": "The Bank_id_code Schema",
          "default": "",
          "examples": [
            "GBDSC"
          ],
          "pattern": "^(.*)$"
        },
        "name": {
          "$id": "#/properties/beneficiary_party/properties/name",
          "type": "string",
          "title": "The Name Schema",
          "default": "",
          "examples": [
            "Wilfred Jeremiah Owens"
          ],
          "pattern": "^(.*)$"
        }
      }
    },
    "charges_information": {
      "$id": "#/properties/charges_information",
      "type": "object",
      "title": "The Charges_information Schema",
      "required": [
        "bearer_code",
        "sender_charges",
        "receiver_charges_amount",
        "receiver_charges_currency"
      ],
      "properties": {
        "bearer_code": {
          "$id": "#/properties/charges_information/properties/bearer_code",
          "type": "string",
          "title": "The Bearer_code Schema",
          "default": "",
          "examples": [
            "SHAR"
          ],
          "pattern": "^(.*)$"
        },
        "sender_charges": {
          "$id": "#/properties/charges_information/properties/sender_charges",
          "type": "array",
          "title": "The Sender_charges Schema",
          "items": {
            "$id": "#/properties/charges_information/properties/sender_charges/items",
            "type": "object",
            "title": "The Items Schema",
            "required": [
              "amount",
              "currency"
            ],
            "properties": {
              "amount": {
                "$id": "#/properties/charges_information/properties/sender_charges/items/properties/amount",
                "type": "string",
                "title": "The Amount Schema",
                "default": "",
                "examples": [
                  "5.00"
                ],
                "pattern": "^(.*)$"
              },
              "currency": {
                "$id": "#/properties/charges_information/properties/sender_charges/items/properties/currency",
                "type": "string",
                "title": "The Currency Schema",
                "default": "",
                "examples": [
                  "GBP"
                ],
                "pattern": "^(.*)$"
              }
            }
          }
        },
        "receiver_charges_amount": {
          "$id": "#/properties/charges_information/properties/receiver_charges_amount",
          "type": "string",
          "title": "The Receiver_charges_amount Schema",
          "default": "",
          "examples": [
            "1.00"
          ],
          "pattern": "^(.*)$"
        },
        "receiver_charges_currency": {
          "$id": "#/properties/charges_information/properties/receiver_charges_currency",
          "type": "string",
          "title": "The Receiver_charges_currency Schema",
          "default": "",
          "examples": [
            "USD"
          ],
          "pattern": "^(.*)$"
        }
      }
    },
    "currency": {
      "$id": "#/properties/currency",
      "type": "string",
      "title": "The Currency Schema",
      "default": "",
      "examples": [
        "GBP"
      ],
      "pattern": "^(.*)$"
    },
    "debtor_party": {
      "$id": "#/properties/debtor_party",
      "type": "object",
      "title": "The Debtor_party Schema",
      "required": [
        "account_name",
        "account_number",
        "account_number_code",
        "address",
        "bank_id",
        "bank_id_code",
        "name"
      ],
      "properties": {
        "account_name": {
          "$id": "#/properties/debtor_party/properties/account_name",
          "type": "string",
          "title": "The Account_name Schema",
          "default": "",
          "examples": [
            "EJ Brown Black"
          ],
          "pattern": "^(.*)$"
        },
        "account_number": {
          "$id": "#/properties/debtor_party/properties/account_number",
          "type": "string",
          "title": "The Account_number Schema",
          "default": "",
          "examples": [
            "GB29XABC10161234567801"
          ],
          "pattern": "^(.*)$"
        },
        "account_number_code": {
          "$id": "#/properties/debtor_party/properties/account_number_code",
          "type": "string",
          "title": "The Account_number_code Schema",
          "default": "",
          "examples": [
            "IBAN"
          ],
          "pattern": "^(.*)$"
        },
        "address": {
          "$id": "#/properties/debtor_party/properties/address",
          "type": "string",
          "title": "The Address Schema",
          "default": "",
          "examples": [
            "10 Debtor Crescent Sourcetown NE1"
          ],
          "pattern": "^(.*)$"
        },
        "bank_id": {
          "$id": "#/properties/debtor_party/properties/bank_id",
          "type": "string",
          "title": "The Bank_id Schema",
          "default": "",
          "examples": [
            "203301"
          ],
          "pattern": "^(.*)$"
        },
        "bank_id_code": {
          "$id": "#/properties/debtor_party/properties/bank_id_code",
          "type": "string",
          "title": "The Bank_id_code Schema",
          "default": "",
          "examples": [
            "GBDSC"
          ],
          "pattern": "^(.*)$"
        },
        "name": {
          "$id": "#/properties/debtor_party/properties/name",
          "type": "string",
          "title": "The Name Schema",
          "default": "",
          "examples": [
            "Emelia Jane Brown"
          ],
          "pattern": "^(.*)$"
        }
      }
    },
    "end_to_end_reference": {
      "$id": "#/properties/end_to_end_reference",
      "type": "string",
      "title": "The End_to_end_reference Schema",
      "default": "",
      "examples": [
        "Wil piano Jan"
      ],
      "pattern": "^(.*)$"
    },
    "fx": {
      "$id": "#/properties/fx",
      "type": "object",
      "title": "The Fx Schema",
      "required": [
        "contract_reference",
        "exchange_rate",
        "original_amount",
        "original_currency"
      ],
      "properties": {
        "contract_reference": {
          "$id": "#/properties/fx/properties/contract_reference",
          "type": "string",
          "title": "The Contract_reference Schema",
          "default": "",
          "examples": [
            "FX123"
          ],
          "pattern": "^(.*)$"
        },
        "exchange_rate": {
          "$id": "#/properties/fx/properties/exchange_rate",
          "type": "string",
          "title": "The Exchange_rate Schema",
          "default": "",
          "examples": [
            "2.00000"
          ],
          "pattern": "^(.*)$"
        },
        "original_amount": {
          "$id": "#/properties/fx/properties/original_amount",
          "type": "string",
          "title": "The Original_amount Schema",
          "default": "",
          "examples": [
            "200.42"
          ],
          "pattern": "^(.*)$"
        },
        "original_currency": {
          "$id": "#/properties/fx/properties/original_currency",
          "type": "string",
          "title": "The Original_currency Schema",
          "default": "",
          "examples": [
            "USD"
          ],
          "pattern": "^(.*)$"
        }
      }
    },
    "numeric_reference": {
      "$id": "#/properties/numeric_reference",
      "type": "string",
      "title": "The Numeric_reference Schema",
      "default": "",
      "examples": [
        "1002001"
      ],
      "pattern": "^(.*)$"
    },
    "payment_id": {
      "$id": "#/properties/payment_id",
      "type": "string",
      "title": "The Payment_id Schema",
      "default": "",
      "examples": [
        "123456789012345678"
      ],
      "pattern": "^(.*)$"
    },
    "payment_purpose": {
      "$id": "#/properties/payment_purpose",
      "type": "string",
      "title": "The Payment_purpose Schema",
      "default": "",
      "examples": [
        "Paying for goods/services"
      ],
      "pattern": "^(.*)$"
    },
    "payment_scheme": {
      "$id": "#/properties/payment_scheme",
      "type": "string",
      "title": "The Payment_scheme Schema",
      "default": "",
      "examples": [
        "FPS"
      ],
      "pattern": "^(.*)$"
    },
    "payment_type": {
      "$id": "#/properties/payment_type",
      "type": "string",
      "title": "The Payment_type Schema",
      "default": "",
      "examples": [
        "Credit"
      ],
      "pattern": "^(.*)$"
    },
    "processing_date": {
      "$id": "#/properties/processing_date",
      "type": "string",
      "title": "The Processing_date Schema",
      "default": "",
      "examples": [
        "2017-01-18"
      ],
      "pattern": "^(.*)$"
    },
    "reference": {
      "$id": "#/properties/reference",
      "type": "string",
      "title": "The Reference Schema",
      "default": "",
      "examples": [
        "Payment for Em's piano lessons"
      ],
      "pattern": "^(.*)$"
    },
    "scheme_payment_sub_type": {
      "$id": "#/properties/scheme_payment_sub_type",
      "type": "string",
      "title": "The Scheme_payment_sub_type Schema",
      "default": "",
      "examples": [
        "InternetBanking"
      ],
      "pattern": "^(.*)$"
    },
    "scheme_payment_type": {
      "$id": "#/properties/scheme_payment_type",
      "type": "string",
      "title": "The Scheme_payment_type Schema",
      "default": "",
      "examples": [
        "ImmediatePayment"
      ],
      "pattern": "^(.*)$"
    },
    "sponsor_party": {
      "$id": "#/properties/sponsor_party",
      "type": "object",
      "title": "The Sponsor_party Schema",
      "required": [
        "account_number",
        "bank_id",
        "bank_id_code"
      ],
      "properties": {
        "account_number": {
          "$id": "#/properties/sponsor_party/properties/account_number",
          "type": "string",
          "title": "The Account_number Schema",
          "default": "",
          "examples": [
            "56781234"
          ],
          "pattern": "^(.*)$"
        },
        "bank_id": {
          "$id": "#/properties/sponsor_party/properties/bank_id",
          "type": "string",
          "title": "The Bank_id Schema",
          "default": "",
          "examples": [
            "123123"
          ],
          "pattern": "^(.*)$"
        },
        "bank_id_code": {
          "$id": "#/properties/sponsor_party/properties/bank_id_code",
          "type": "string",
          "title": "The Bank_id_code Schema",
          "default": "",
          "examples": [
            "GBDSC"
          ],
          "pattern": "^(.*)$"
        }
      }
    }
  }
}`
