---
title: "Why reporting rate beats click rate in phishing simulations"
slug: "reporting-rate-beats-click-rate"
date: 2026-08-22T10:00:00+03:00
lastmod: 2026-08-22T10:00:00+03:00
author: "spectrumsec"
summary: "Click rate is the metric everyone reports and the one that matters least. Here is the case for optimising how fast your people report instead."
tags: ["phishing", "awareness", "metrics"]
draft: false
toc: true
---

<!-- TODO: placeholder article. Replace with real content and data. -->

Most phishing-simulation reports lead with **click rate** — the percentage of recipients
who clicked the link. It is easy to measure and easy to put on a slide. It is also the
wrong thing to optimise.

## The problem with click rate

A determined attacker only needs one click. Driving your click rate from 12% to 8% does
not meaningfully change that. What changes your risk is **how quickly someone raises the
alarm**, because that determines how long the attacker operates undetected.

## What to measure instead

1. **Report rate** — what fraction of recipients reported the message.
2. **Time to first report** — minutes between delivery and the first alert.
3. **Report-to-click ratio** — are reporters outpacing clickers over time?

## Making it actionable

- Put a one-click report button in the mail client and make sure it works.
- Thank every reporter, publicly and quickly. Never punish a clicker.
- Track the trend across campaign cycles, not the number from any single run.

If you want help designing a programme around these metrics, our
[phishing simulation](/services/phishing-simulation/) service is built around exactly this.
