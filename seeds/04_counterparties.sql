-- ExchangeOS seeds — 04: counterparties (BIC records for initial CLS members + Brazilian banks).
-- Populates bic_records globally. Tenant-scoped counterparty rows (table `counterparties`
-- created in 000002) are intentionally NOT seeded here — those are created at tenant onboarding.

BEGIN;

INSERT INTO bic_records (bic, institution_name, country, lei, active) VALUES
    -- ── CLS settlement member banks (sample) ────────────────────────────────
    ('CLSBUS33', 'CLS Bank International',           'US', '254900V4Q4YE3WJM7Q23', true),
    ('DEUTDEFF', 'Deutsche Bank AG',                  'DE', '7LTWFZYICNSX8D621K86', true),
    ('CHASUS33', 'JPMorgan Chase Bank N.A.',          'US', '7H6GLXDRUGQFU57RNE97', true),
    ('CITIUS33', 'Citibank N.A.',                     'US', 'E57ODZWZ7FF32TWEFA76', true),
    ('BOFAUS3N', 'Bank of America N.A.',              'US', 'B4TYDEB6GKMZO031MB27', true),
    ('GSCMUS33', 'Goldman Sachs Bank USA',            'US', '784F5XWPLTWKTBV3E584', true),
    ('MSNYUS33', 'Morgan Stanley & Co. LLC',          'US', '9R7GPTSO7KV3UQJZQ078', true),
    ('UBSWCHZH', 'UBS Switzerland AG',                'CH', '549300WOIFUSNYH0FL22', true),
    ('CRESCHZZ', 'Credit Suisse (Schweiz) AG',        'CH', '549300CWO1BYAVHFGI82', true),
    ('BARCGB22', 'Barclays Bank UK PLC',              'GB', '54930032RWXBO1OOL004', true),
    ('HSBCGB2L', 'HSBC Bank plc',                     'GB', 'MP6I5ZYZBEU3UXPYFY54', true),
    ('NWBKGB2L', 'NatWest Markets plc',               'GB', 'RR3QWICWWIPCS8A4S074', true),
    ('BNPAFRPP', 'BNP Paribas',                       'FR', 'R0MUWSFPU8MPRO8K5P83', true),
    ('SOGEFRPP', 'Société Générale',                  'FR', 'O2RNE8IBXP4R0TD8PU41', true),
    ('AGRIFRPP', 'Crédit Agricole CIB',               'FR', '1VUV7VQFKUOQSJ21A208', true),
    ('SMBCJPJT', 'Sumitomo Mitsui Banking Corp.',     'JP', '353800F562J3UTPV6E63', true),
    ('BOTKJPJT', 'MUFG Bank',                         'JP', '7VOTYQBDDXOK2GDJJ7Q3', true),
    ('MHCBJPJT', 'Mizuho Bank',                       'JP', 'RB0PEZSDGCO3JS6CEU02', true),
    ('NABAAU3M', 'National Australia Bank',           'AU', 'F8SB4JFBSYQFRQEH3Z21', true),
    ('CTBAAU2S', 'Commonwealth Bank of Australia',    'AU', 'MSFSBD3QN1GSN7Q6C537', true),
    ('ANZBAU3M', 'Australia and New Zealand Banking Group','AU','JHE42UYNWWTJB8YTTU19',true),
    ('WPACAU2F', 'Westpac Banking Corporation',       'AU', 'EN5TNI6CI43VEPAMHL14', true),
    ('RBCBCATT', 'Royal Bank of Canada',              'CA', 'ES7IP3U3RHIGC71XBU11', true),
    ('TDOMCATTTOR','Toronto-Dominion Bank',           'CA', 'PT3QB789TSUIDF371261', true),
    ('NORDDKKK', 'Nordea Bank',                       'DK', '529900ODI3047E2LIV03', true),

    -- ── Brazilian banks (BRL counterparties) ────────────────────────────────
    ('BCBRBRBR', 'Banco Central do Brasil',           'BR', '529900P2N1KQK8WHN398', true),
    ('ITAUBRSP', 'Itaú Unibanco S.A.',                'BR', '95KGCWYIVAFTYI5OTI23', true),
    ('BBRSBRRJ', 'Banco do Brasil S.A.',              'BR', '5493005JG4U3BIDFM812', true),
    ('CAIXBRRJ', 'Caixa Econômica Federal',           'BR', '529900FNUW3P8Q7EIJ02', true),
    ('BRADBRRJ', 'Banco Bradesco S.A.',               'BR', '529900T7BAS2HCH9KX95', true),
    ('SANBBRSP', 'Banco Santander (Brasil) S.A.',     'BR', '549300I8U5NXOA6OQ737', true),
    ('BCITITMM', 'Banco BTG Pactual S.A.',            'BR', '5493004QVZJZ7VRKEU38', true),
    ('XPMTBRSP', 'XP Investimentos CCTVM S.A.',       'BR', '549300S5UU7RV80LBP06', true),

    -- ── CFETS (China onshore reference) ─────────────────────────────────────
    ('ICBKCNBJ', 'Industrial and Commercial Bank of China','CN','549300DTUYXVMJXZNY75',true),
    ('BKCHCNBJ', 'Bank of China',                     'CN', '549300VGCB3MTAREV382', true),
    -- CFETS has no official SWIFT BIC; this is an illustrative 8-char ISO 9362
    -- form (CFET=bank, CN=country, SH=location). Replace with the real BIC
    -- once CFETS publishes one via correspondent banking.
    ('CFETCNSH','China Foreign Exchange Trade System','CN','5493002Q6N5J3D6FUH26',true)
ON CONFLICT (bic) DO NOTHING;

COMMIT;
