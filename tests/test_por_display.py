#!/usr/bin/env python3
"""Unit checks for skills/lib/por_display.py human-facing reserve formatting."""

import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "skills" / "lib"))

import por_display  # noqa: E402


class FormatUsd(unittest.TestCase):
    def test_onchain_zero_is_unavailable(self):
        self.assertEqual(por_display.format_usd_onchain(0), "unavailable")

    def test_large_usd_uses_commas(self):
        self.assertEqual(
            por_display.format_usd_onchain(132_692_864_216),
            "$132,692,864,216 USD (on-chain storage is whole USD)",
        )
        self.assertEqual(
            por_display.format_usd_float(10_645_747_991.82),
            "$10,645,747,991.82 USD",
        )

    def test_sentinel_one_prefers_assessment_float(self):
        # Encode rule: positive sub-dollar totals round to on-chain 1.
        self.assertEqual(
            por_display.format_usd_onchain(1, assessment_float=0.42),
            "$0.42 USD",
        )
        self.assertEqual(
            por_display.format_usd_onchain(1),
            "$1 USD (on-chain storage is whole USD)",
        )


class FormatRatioBps(unittest.TestCase):
    def test_unavailable_sentinel(self):
        self.assertEqual(por_display.format_reserve_ratio_bps(65535), "unavailable")
        self.assertEqual(por_display.format_reserve_ratio_bps(0), "0.00%")

    def test_one_hundred_percent(self):
        self.assertEqual(por_display.format_reserve_ratio_bps(10000), "100.00%")

    def test_two_decimal_bps(self):
        self.assertEqual(por_display.format_reserve_ratio_bps(10003), "100.03%")

    def test_saturated_cap(self):
        self.assertEqual(por_display.format_reserve_ratio_bps(65534), "≥655.34%")
        self.assertEqual(por_display.format_reserve_ratio_rate(13.47805109), "≥655.34%")


class BuildDisplay(unittest.TestCase):
    def test_lines_label_self_reported(self):
        display = por_display.build_reserve_display(
            {
                "schemaVersion": 7,
                "reserveTotalsMeaningful": True,
                "totalLiabilityUsd": 1_000_000,
                "totalReserveUsd": 0,
                "btcReserveRatioBps": 10000,
                "ethReserveRatioBps": 65535,
                "usdtReserveRatioBps": 10003,
                "usdcReserveRatioBps": 65534,
            },
            {
                "totalLiabilityUsd": 1_000_000.4,
                "totalReserveUsd": 0.0,
                "majorCoinReserveRates": [
                    {"coin": "BTC", "reserveRate": 1.0, "selfReported": True},
                ],
            },
        )
        self.assertIn("Total liabilities: $1,000,000.40 USD (self-reported)", display["lines"])
        self.assertIn("Total reserves: $0.00 USD (self-reported)", display["lines"])
        self.assertIn("BTC reserve ratio: 100.00% (self-reported)", display["lines"])
        self.assertIn("ETH reserve ratio: unavailable (self-reported)", display["lines"])
        self.assertIn("USDT reserve ratio: 100.03% (self-reported)", display["lines"])
        self.assertIn("USDC reserve ratio: ≥655.34% (self-reported)", display["lines"])
        usdt = display["majorCoinReserveRatios"][2]
        self.assertEqual(usdt["coin"], "USDT")
        self.assertEqual(usdt["rawBps"], 10003)
        self.assertIn("not a solvency opinion", display["note"].lower())

    def test_prev7_zeros_are_unavailable(self):
        display = por_display.build_reserve_display(
            {
                "schemaVersion": 6,
                "reserveTotalsMeaningful": False,
                "totalLiabilityUsd": 0,
                "totalReserveUsd": 0,
                "btcReserveRatioBps": 65535,
                "ethReserveRatioBps": 65535,
                "usdtReserveRatioBps": 65535,
                "usdcReserveRatioBps": 65535,
            },
        )
        self.assertEqual(display["totalLiabilityUsd"]["display"], "unavailable")
        self.assertEqual(display["majorCoinReserveRatios"][0]["display"], "unavailable")


if __name__ == "__main__":
    unittest.main()
