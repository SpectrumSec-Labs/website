---
title: "Board Briefing — September 2026: patched is not the same as safe"
slug: "board-briefing-2026-09"
date: 2026-10-10T10:00:00+03:00
author: "spectrumsec"
summary: "Remote-access gateways were exploited for weeks before a fix existed. Plus: new EU reporting duties with a 24-hour clock, and a record quarter for ransomware."
tags: ["board-briefing", "remote-access", "cra", "ai-act", "ransomware"]
categories: ["Board Briefing"]
draft: false
toc: false
---

## The month in one paragraph

In September, attackers broke into the devices that many organisations use to let staff work
remotely — and they did it for weeks before the manufacturer knew there was a problem, let
alone had a fix. When the patch finally arrived, the official advice was not "install it and
move on" but "check whether you were already breached first". At the same time, a new EU law
started the clock on a 24-hour reporting duty for companies that make digital products, and
ransomware closed the quarter at a record high. The theme for leadership: speed matters, but
assurance matters more. "We patched" is the start of the answer, not the end of it.

## The story this month: the front door was open before the lock existed

On 27 September, Citrix disclosed critical flaws in NetScaler ADC and NetScaler Gateway — the
equipment that sits at the edge of many corporate networks and gives employees remote access
to internal systems. The US cybersecurity agency CISA warned the same day that two of them
were being actively exploited worldwide, and France's CERT-FR issued an alert of its own.

The detail that matters most came from Google's Mandiant team three days later: the campaign
had been running **since at least early September** — weeks before the public disclosure.
Organisations in government, financial services, technology, education, and legal and
professional services across Europe and North America were likely affected. Once inside,
the attackers planted hidden backdoors on the devices, then used them as a tunnel into
internal networks to look around and steal credentials. Who is behind it has not been
publicly confirmed.

This is not a one-off. NetScaler devices were at the centre of major exploitation waves in
2023 and 2025 as well, and internet-facing gateways from many vendors are a favourite target
for one simple reason: they are exposed to the whole internet, they hold the keys to
everything behind them, and they are often less closely watched than laptops and servers.

What the authorities told organisations to do is the real lesson:

- **Check for compromise before patching** — and preserve evidence, because updating can
  erase the traces an investigator needs.
- **Assume credentials were stolen.** After patching, rotate passwords, keys and
  certificates, and cut off existing sessions.
- **Get logs off the device.** These appliances overwrite their own logs quickly; without a
  copy elsewhere, you may never know what happened.

In other words, an organisation that patched promptly on 27 September could still have been
breached on 10 September — and would only find out if it went looking.

## Regulation and compliance: two EU dates that moved this quarter

**The Cyber Resilience Act's reporting duty is live (11 September).** Manufacturers of
products with digital elements sold in the EU — software, connected devices, hardware with
software inside — must now report actively exploited vulnerabilities and severe security
incidents affecting those products. The clock is tight: an early warning within **24 hours**
of becoming aware, a fuller notification within 72 hours, and a final report within 14 days
of a fix (or a month, for incidents). Reports go through a single EU platform run by ENISA,
which opened the same day. Legal commentators note that the duty covers products already on
the market, not just new ones. Fines for breaching it reach €15 million or 2.5% of global
turnover. The NetScaler story shows exactly the situation the law was written for — and if
your company makes or sells digital products, it now applies to you.

**The AI Act's high-risk rules were delayed — but not everything was.** An amending regulation
in force since 27 July moved the obligations for "high-risk" AI systems (for example in
hiring, credit scoring or education) to 2 December 2027, and for AI built into regulated
products to 2 August 2028. Transparency duties were not postponed: since 2 August 2026,
people must be told when they are dealing with an AI system, and AI-generated content must be
identifiable (systems already on the market have until December to comply with the labelling
part). Delay is not cancellation; the time is best used to inventory where AI is
already in use.

## The wider picture: a record quarter for ransomware

One tracker counted 2,627 ransomware attacks between July and September — the highest
quarter on record, up 29% on the previous quarter and 61% on a year earlier. Only 247 were
confirmed by the victims themselves, so the real picture is murkier than any headline number.
Manufacturing was the most targeted business sector. In Europe, the Berlin state government
was hit and refused a $2.3 million demand. Ransomware is no longer an exceptional event to
plan for; it is the background condition.

## What this means for you

- **Know your front doors.** Keep an up-to-date list of every internet-facing device and who
  is responsible for it — remote-access gateways, firewalls, VPNs, file-transfer servers.
- **Plan for the emergency patch, not just the routine one.** Agree in advance who can take
  a critical system offline at short notice, and how the business keeps running while it is.
- **Pay for visibility, not only for prevention.** Logs stored off the device, and someone
  who looks at them, are what turn "we hope we're fine" into "we know".
- **Budget to replace, not just repair.** Edge equipment that is out of support should be
  on the board's risk register with a date attached.
- **If you make digital products, test your 24-hour CRA process now** — including who has
  access to the EU reporting platform — before you need it for real.

## Three questions to ask your security team this month

1. Which of our systems are directly reachable from the internet, and how quickly could we
   take each one offline if its vendor announced an emergency?
2. When we patch a critical flaw, do we also check whether we were breached *before* the
   patch — and who would have to sign off on that investigation?
3. If we are a manufacturer under the Cyber Resilience Act, who would file our 24-hour early
   warning this weekend, and have they ever practised it?

*If you would like an independent view of your internet-facing exposure, our
[penetration testing](/services/penetration-testing/) and
[infrastructure hardening](/services/infrastructure-hardening/) services start exactly there.*

## Sources

- Citrix NetScaler zero-days: [CISA alert, 27 September 2026 (revised 9
  October)](https://www.cisa.gov/news-events/alerts/2026/09/27/critical-zero-day-vulnerabilities-exploited-citrix-netscaler-adc-gateway)
  · [Google Threat Intelligence / Mandiant, 30 September
  2026](https://cloud.google.com/blog/topics/threat-intelligence/defending-against-active-exploitation-of-citrix-netscaler-adc-and-gateway-appliances)
  · [CERT-FR alert (French)](https://www.cert.ssi.gouv.fr/alerte/CERTFR-2026-ALE-011/)
- Cyber Resilience Act reporting: [European Commission — CRA reporting
  obligations](https://digital-strategy.ec.europa.eu/en/policies/cra-reporting) · [Hunton
  Andrews Kurth — scope and
  penalties](https://www.hunton.com/privacy-and-cybersecurity-law-blog/eu-cyber-resilience-act-reporting-obligations-take-effect-for-manufacturers)
- AI Act amendments: [Regulation (EU) 2026/1744, Official
  Journal](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=OJ:L_202601744) · [Northern
  Ireland Assembly — published EU acts
  tracker](https://www.niassembly.gov.uk/assembly-business/committees/2022-2027/windsor-framework-democratic-scrutiny-committee/eu-acts/published-eu-acts/regeu20261744/)
- Ransomware figures: [Comparitech — Ransomware roundup Q3
  2026](https://www.comparitech.com/news/ransomware-roundup-q3-2026-stats-on-attacks-ransoms-and-active-gangs)
