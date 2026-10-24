---
title: "Board Briefing — October 2026: when your supplier goes down, so do you"
slug: "board-briefing-2026-10"
date: 2026-10-10T12:00:00+03:00
author: "spectrumsec"
summary: "One cloud provider's ransomware outage took 495 organisations offline at once. What that means for your suppliers, your backups and your NIS2 duties."
tags: ["board-briefing", "supply-chain", "cloud", "nis2"]
categories: ["Board Briefing"]
draft: true
toc: false
---

## The month in one paragraph

On 7 October a ransomware attack shut down part of a major Japanese cloud platform, and
495 businesses and local governments that rented capacity there went down with it — not
because they were attacked, but because their supplier was. The provider then told them
that much of their data would be difficult to recover, and that they should restore from
their own backups. Three months earlier, a global pharmaceutical company disclosed that
patient and proprietary data had been stolen from cloud environments run by third parties.
The lesson for leadership is the same in both cases: you can outsource your systems, but
you cannot outsource the consequences. And in the EU, the law now says so explicitly.

## The story this month: one supplier, hundreds of victims

**IDCF Cloud (Japan).** IDC Frontier, a SoftBank subsidiary, confirmed that ransomware forced
it to shut down virtual servers across four zones of one of its cloud regions in the early
hours of 7 October. It named 495 affected customers, including private companies and local
authorities. The most important line in its notice was not about the attack itself: customer
data in the affected zones is expected to be hard to restore, and customers were directed to
recover from backups *they* hold.

**Nissui (Japan).** The same day, seafood group Nissui announced that its logistics
subsidiary could not receive or ship goods, after suspected unauthorised access to an
outsourced data centre. Press reports link the outage to the IDCF Cloud attack; Nissui has
not named its provider. Either way, a food company stopped moving product because of a
system it did not run.

**Amgen (United States).** In July, Amgen told investors that proprietary data and patient
health information had been taken from cloud environments hosted by third-party providers.
Its products and manufacturing were unaffected — but it still classed the incident as
material, and the duty to notify patients and regulators sits with Amgen, not the provider.

Three different failures, one pattern:

- **Concentration.** When many organisations rely on one provider, one attack becomes
  hundreds of outages at once — and every victim queues for the same recovery team.
- **Shared responsibility is not shared liability.** Cloud providers keep the platform
  running. Keeping *your* data recoverable, and answering for it when it leaks, is yours.
- **Operational dependency is invisible until it fails.** Few boards could name, today,
  which suppliers would stop them trading within 24 hours.

## Regulation and compliance: NIS2 puts this on the board's desk

If your organisation is in scope of NIS2 — the EU's cybersecurity directive for essential
and important entities, now law in most member states — supplier risk is not best practice;
it is a named legal requirement:

- **Supply-chain security is mandatory.** Article 21 lists the security of relationships
  with direct suppliers and service providers among the measures every in-scope entity must
  take.
- **Management must approve and oversee it.** Article 20 requires the management body to
  approve those measures, oversee their implementation, and take regular training — and
  makes its members answerable for failures.
- **The sanctions are real.** Maximum fines must be at least €10 million or 2% of global
  turnover, whichever is higher, for essential entities (€7 million or 1.4% for important
  ones), and member states may set them higher. If an essential entity ignores orders to fix
  shortcomings, supervisors can seek a temporary ban on the responsible executive managing it.

Germany shows where this is heading. Its implementing law has been in force since 6 December
2025, brings roughly 29,500 organisations into scope, and requires management to approve the
risk measures, oversee them and attend training — with managers liable to their own company
for damage caused by failing to do so. Supervisors across the EU have so far announced few
fines, but they are moving from registration to audits. The first question in an audit is
likely to be the simplest: *show us that management approved this.*

## What this means for you

- **Map critical suppliers by impact, not by spend.** Your largest contract is not
  necessarily your largest risk. Rank suppliers by how fast their failure would stop you
  trading.
- **Own a copy of what you cannot lose.** Critical data should have a backup that sits
  outside your main provider and that you have actually restored from recently.
- **Put supplier failure into your incident plans.** Most response plans assume *we* are
  attacked. Rehearse the scenario where a supplier is — and you learn about it from the news.
- **Check your contracts.** How quickly must a supplier tell you about an incident? Can you
  get your data out, in a usable form, if they cannot recover?
- **Minute it.** If you are in scope of NIS2, board approval of supplier-security measures
  should be recorded, current and easy to produce.

## Three questions to ask your security team this month

1. Which three suppliers, if they went dark tomorrow, would stop us trading — and how long
   could we keep operating without each of them?
2. If our main cloud provider told us our data was unrecoverable, where is the copy we
   control, and when did we last prove we can restore from it?
3. Has the board formally approved our supplier-security requirements this year, and could
   we show a regulator the evidence?

*If you would like an independent view of your supplier exposure, our
[security assessment](/services/security-assessment/) maps critical dependencies and tests
how quickly you could recover without them.*

## Sources

- BleepingComputer — [Ransomware attack disrupts Japan's IDCF Cloud used by govt
  clients](https://bleepingcomputer.com/news/security/ransomware-attack-disrupts-japans-idcf-cloud-used-by-govt-clients)
- ASCII.jp — [IDCF Frontier: customer data recovery challenging after ransomware attack,
  affecting 495 entities](https://ascii.jp/en/elem/000/004/441/4441355/)
- Security Measures Lab (Japanese) — [Nissui Logistics system outage and shipment
  stoppage](https://rocket-boys.co.jp/security-measures-lab/nissui-logistics-system-outage-cyberattack-2026)
- Amgen Inc. — [Form 8-K, cybersecurity incident disclosure (SEC
  EDGAR)](https://www.sec.gov/Archives/edgar/data/0000318154/000031815426000119/amgn-20260729.htm)
- EUR-Lex — [Directive (EU) 2022/2555 (NIS2), Articles 20, 21, 32 and
  34](https://eur-lex.europa.eu/eli/dir/2022/2555/oj)
- YPOG — [Germany's NIS2 Implementation Act](https://www.ypog.law/en/insight/germanys-nis2-implementation-act)
