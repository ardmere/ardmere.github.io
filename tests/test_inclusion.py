"""Tests for the Ardmere PoR skills.

Two suites:
  * Inclusion proof schemes — parity with the expected roots the Go tests in
    internal/exchanges/*/merkle assert, so a pass here means the Python port and
    the shipped Go implementation agree.
  * Custody risk score — the decision table in
    skills/por-custody-risk/references/decision-table.md, plus the anchor-read
    short-circuits, run against recorded anchor fixtures.

Run: python3 -m unittest discover -s tests -v
"""

import json
import re
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
INCLUSION_DIR = ROOT / "skills" / "por-inclusion-proof"
CUSTODY_DIR = ROOT / "skills" / "por-custody-risk"
INCLUSION = INCLUSION_DIR / "scripts" / "inclusion_check.py"
CUSTODY = CUSTODY_DIR / "scripts" / "risk_score.py"
FIXTURES = ROOT / "fixtures"
ANCHOR_FIXTURE = FIXTURES / "anchor" / "sepolia-v3-records.json"

sys.path.insert(0, str(ROOT / "skills" / "lib"))


def run(script, *args, cwd=ROOT):
    return subprocess.run([sys.executable, str(script), *args],
                          capture_output=True, text=True, cwd=cwd)


# exchange, fixture, expected root asserted by the Go test of the same scheme.
CASES = [
    ("bitget", "bitget/merkel_tree_bg.json", "ca89456bb711c913"),
    ("bingx", "bingx/mock_user_proof.json",
     "e12dba68a58da352b4fd8901f3d989e361cba37ce1fc4b28319b285c42244347"),
    ("bitunix", "bitunix/merkle_tree_bitunix.json",
     "7d409a80ebbbf4cef17f3f9c5f56f6e26040f1ca285395897cf695b21391679d"),
    ("bybit", "bybit/mock_user_merkle_tree_path_40_v5.json",
     "c91d5de0554f97244e4d9f8056fad70fa0cb2cdb23c290ec597a042645dcbc03"),
    ("coinex", "coinex/mock_user_proof.json",
     "8d44e525e0ae827da3d27f8f27b50d4c5bc8d94b03b7a29bbf7219a8ba9cf05f"),
    ("mexc", "mexc/merkle_proof.json", "35d01ffef2beec3ff0166c03f3d9b0cc"),
    ("ourbit", "ourbit/merkel_tree_ourbit.json", "df1d34512abd50712b9f5713b51fad8e"),
    ("toobit", "toobit/user_1_demo_proof.json",
     "3586fdd20930df70e214ca0157c68bb6226289824f05f92d5e0465eb5297f727"),
]

TRUNCATED = {"bitget", "mexc", "ourbit"}

# Schemes whose assessment in this repo publishes a merkle root of the same width.
BOUND = {
    "bybit": "docs/reports/bybit/2025061709-assessment.json",
    "bitunix": "docs/reports/bitunix/20261103-assessment.json",
    "coinex": "docs/reports/coinex/20260914-assessment.json",
    "toobit": "docs/reports/toobit/20261001-assessment.json",
}


class SchemeList(unittest.TestCase):
    def test_list_schemes_is_exactly_what_is_implemented(self):
        proc = run(INCLUSION, "--list-schemes")
        self.assertEqual(proc.returncode, 0, proc.stderr)
        listed = [line for line in proc.stdout.split() if line]
        self.assertEqual(sorted(listed), sorted(exchange for exchange, _f, _r in CASES))

    def test_list_schemes_matches_the_go_merkle_packages(self):
        """Every listed scheme has a Go package this port was checked against."""
        listed = set(run(INCLUSION, "--list-schemes").stdout.split())
        for exchange in listed:
            self.assertTrue((ROOT / "internal" / "exchanges" / exchange / "merkle").is_dir(),
                            f"{exchange} has no Go merkle package")


class ParityWithGo(unittest.TestCase):
    def test_every_fixture_recomputes_the_go_test_root(self):
        for exchange, fixture, expected in CASES:
            with self.subTest(exchange=exchange):
                path = FIXTURES / fixture
                self.assertTrue(path.is_file(), f"{path} is missing")
                proc = run(INCLUSION, "--exchange", exchange, str(path))
                self.assertEqual(proc.returncode, 0, proc.stderr)
                doc = json.loads(proc.stdout)
                self.assertEqual(doc["verdict"], "PASS")
                self.assertEqual(doc["recomputedRoot"], expected)

    def test_the_go_verifiers_agree_on_the_same_fixtures(self):
        """Optional cross-check: run the Go tests when a toolchain is present.

        The Python suite already pins the roots, so this is a second opinion rather
        than the contract. It is skipped where `go` is absent, which keeps the skill
        usable from a bare checkout.
        """
        go = shutil.which("go")
        if not go:
            self.skipTest("go is not installed; the pinned roots above are the check")
        proc = subprocess.run(
            [go, "test", "./internal/exchanges/bitget/merkle", "./internal/exchanges/bingx/merkle",
             "./internal/exchanges/bitunix/merkle", "./internal/exchanges/bybit/merkle",
             "./internal/exchanges/coinex/merkle", "./internal/exchanges/mexc/merkle",
             "./internal/exchanges/ourbit/merkle", "./internal/exchanges/toobit/merkle"],
            capture_output=True, text=True, cwd=ROOT, timeout=300,
        )
        self.assertEqual(proc.returncode, 0, proc.stdout + proc.stderr)

    def test_truncated_schemes_say_so(self):
        for exchange, fixture, _root in CASES:
            with self.subTest(exchange=exchange):
                doc = json.loads(run(INCLUSION, "--exchange", exchange, str(FIXTURES / fixture)).stdout)
                ids = [f["id"] for f in doc["findings"]]
                if exchange in TRUNCATED:
                    self.assertIn("TRUNCATED_HASH", ids)
                    self.assertLess(doc["rootWidthBytes"] * 8, 256)
                else:
                    self.assertNotIn("TRUNCATED_HASH", ids)
                    self.assertEqual(doc["rootWidthBytes"] * 8, 256)

    def test_every_verdict_carries_its_scope(self):
        for exchange, fixture, _root in CASES:
            with self.subTest(exchange=exchange):
                doc = json.loads(run(INCLUSION, "--exchange", exchange, str(FIXTURES / fixture)).stdout)
                self.assertIn("does not prove the exchange holds enough reserves", doc["scope"])

    def test_no_scheme_invents_a_decision(self):
        for exchange, fixture, _root in CASES:
            with self.subTest(exchange=exchange):
                doc = json.loads(run(INCLUSION, "--exchange", exchange, str(FIXTURES / fixture)).stdout)
                self.assertNotIn("decision", doc)


class EverySchemeRejectsTampering(unittest.TestCase):
    """Each scheme must fail when the file's own declared root is changed.

    The mutation is per-scheme because each file carries the root in a different
    place; what is shared is the expectation: a changed root is a FAIL, not a PASS.
    """

    def setUp(self):
        self.tmp = Path(tempfile.mkdtemp(prefix="por-tamper-"))

    def tearDown(self):
        shutil.rmtree(self.tmp, ignore_errors=True)

    def _mutate_root(self, exchange, proof):
        """Flip one hex character of whatever root the scheme compares against."""
        def flip(text):
            text = str(text)
            head, last = text[:-1], text[-1]
            return head + ("0" if last != "0" else "1")

        if exchange == "bitget":
            proof["root"]["hash"] = flip(proof["root"]["hash"])
        elif exchange == "bingx":
            proof["merkleRoot"] = flip(proof["merkleRoot"])
        elif exchange == "bitunix":
            proof["path"][-1]["merkleLeaf"] = flip(proof["path"][-1]["merkleLeaf"])
        elif exchange == "bybit":
            proof["path"][-1]["hash"] = flip(proof["path"][-1]["hash"])
        elif exchange == "coinex":
            proof["root"]["hash"] = flip(proof["root"]["hash"])
        elif exchange == "mexc":
            proof["pathList"][-1]["hash"] = flip(proof["pathList"][-1]["hash"])
        elif exchange == "ourbit":
            proof["path"][-1]["hash"] = flip(proof["path"][-1]["hash"])
        elif exchange == "toobit":
            proof["root"]["hash"] = flip(proof["root"]["hash"])
        else:
            raise AssertionError(f"no root mutation defined for {exchange}")
        return proof

    def test_every_scheme_fails_on_a_changed_root(self):
        for exchange, fixture, _root in CASES:
            with self.subTest(exchange=exchange):
                proof = self._mutate_root(exchange, json.loads((FIXTURES / fixture).read_text()))
                path = self.tmp / f"{exchange}.json"
                path.write_text(json.dumps(proof))
                proc = run(INCLUSION, "--exchange", exchange, str(path))
                self.assertEqual(proc.returncode, 1, f"{exchange} did not fail: {proc.stdout}")
                self.assertEqual(json.loads(proc.stdout)["verdict"], "FAIL")

    def test_every_scheme_refuses_an_empty_object(self):
        for exchange, _fixture, _root in CASES:
            with self.subTest(exchange=exchange):
                path = self.tmp / "empty.json"
                path.write_text("{}")
                proc = run(INCLUSION, "--exchange", exchange, str(path))
                self.assertIn(proc.returncode, (1,), proc.stdout)
                self.assertTrue(proc.stderr.strip(), f"{exchange} printed no reason")


class NegativeCases(unittest.TestCase):
    def setUp(self):
        self.tmp = Path(tempfile.mkdtemp(prefix="por-tests-"))

    def tearDown(self):
        shutil.rmtree(self.tmp, ignore_errors=True)

    def write(self, name, doc):
        path = self.tmp / name
        path.write_text(json.dumps(doc) if not isinstance(doc, str) else doc)
        return path

    def test_tampered_bingx_leaf_fails(self):
        proof = json.loads((FIXTURES / "bingx/mock_user_proof.json").read_text())
        leaf = proof["merkleLeaf"]
        proof["merkleLeaf"] = leaf[:-1] + ("0" if leaf[-1] != "0" else "1")
        proc = run(INCLUSION, "--exchange", "bingx", str(self.write("t-bingx.json", proof)))
        self.assertEqual(proc.returncode, 1)
        self.assertEqual(json.loads(proc.stdout)["verdict"], "FAIL")
        self.assertIn("PATH_MISMATCH", [f["id"] for f in json.loads(proc.stdout)["findings"]])

    def test_tampered_bitget_balance_fails(self):
        proof = json.loads((FIXTURES / "bitget/merkel_tree_bg.json").read_text())
        asset = sorted(proof["balances"])[0]
        proof["balances"][asset] = str(float(proof["balances"][asset]) + 1)
        proc = run(INCLUSION, "--exchange", "bitget", str(self.write("t-bitget.json", proof)))
        self.assertEqual(proc.returncode, 1)
        self.assertEqual(json.loads(proc.stdout)["verdict"], "FAIL")

    def test_tampered_mexc_account_id_fails_as_self_hash(self):
        proof = json.loads((FIXTURES / "mexc/merkle_proof.json").read_text())
        proof["self"]["accountId"] = str(proof["self"]["accountId"]) + "0"
        proc = run(INCLUSION, "--exchange", "mexc", str(self.write("t-mexc.json", proof)))
        self.assertEqual(proc.returncode, 1)
        doc = json.loads(proc.stdout)
        self.assertEqual(doc["verdict"], "FAIL")
        self.assertIn("SELF_HASH_MISMATCH", [f["id"] for f in doc["findings"]])

    def test_tampered_bybit_balance_fails(self):
        proof = json.loads((FIXTURES / "bybit/mock_user_merkle_tree_path_40_v5.json").read_text())
        proof["self"]["balance"]["BTC"] = "1.00000000"
        proc = run(INCLUSION, "--exchange", "bybit", str(self.write("t-bybit.json", proof)))
        self.assertEqual(proc.returncode, 1)

    def test_absent_root_fails_rather_than_passing(self):
        path = self.write("no-root.json", {"merkleLeaf": "ab" * 32, "merklePath": "", "merkleRoot": ""})
        proc = run(INCLUSION, "--exchange", "bingx", str(path))
        self.assertEqual(proc.returncode, 1)
        self.assertIn("merkleRoot", proc.stderr)

    def test_unknown_exchange_is_refused(self):
        proc = run(INCLUSION, "--exchange", "kraken", str(self.write("unknown.json", {})))
        self.assertEqual(proc.returncode, 2)
        self.assertIn("unsupported exchange", proc.stderr)

    def test_missing_file_is_refused(self):
        proc = run(INCLUSION, "--exchange", "bingx", str(self.tmp / "absent.json"))
        self.assertEqual(proc.returncode, 1)

    def test_wrong_root_width_is_a_width_problem_not_a_failure(self):
        proof = FIXTURES / "bitget/merkel_tree_bg.json"
        assessment = self.write("bitget-wide.json", {
            "exchange": "bitget",
            "snapshot": {"id": "20261021",
                         "exchangeMerkleRoot": "00" * 24 + "1ef42111598bd71b"},
        })
        proc = run(INCLUSION, "--exchange", "bitget", "--assessment", str(assessment), str(proof))
        doc = json.loads(proc.stdout)
        self.assertEqual(doc["verdict"], "PASS")
        self.assertTrue(doc["binding"]["widthsDiffer"])
        self.assertIsNone(doc["binding"]["matches"])
        self.assertIn("ROOT_WIDTH_MISMATCH", [f["id"] for f in doc["findings"]])

    def test_mexc_root_width_mismatch_is_reported_the_same_way(self):
        assessment = self.write("mexc-wide.json", {
            "exchange": "mexc",
            "snapshot": {"id": "202610", "exchangeMerkleRoot": "ab" * 32},
        })
        proc = run(INCLUSION, "--exchange", "mexc", "--assessment", str(assessment),
                   str(FIXTURES / "mexc/merkle_proof.json"))
        doc = json.loads(proc.stdout)
        self.assertEqual(doc["verdict"], "PASS")
        self.assertIsNone(doc["binding"]["matches"])
        self.assertIn("ROOT_WIDTH_MISMATCH", [f["id"] for f in doc["findings"]])

    def test_root_mismatch_with_assessment_fails(self):
        assessment = self.write("bitget-assessment.json", {
            "exchange": "bitget",
            "snapshot": {"id": "202605", "exchangeMerkleRoot": "ab" * 8},
        })
        proc = run(INCLUSION, "--exchange", "bitget", "--assessment", str(assessment),
                   str(FIXTURES / "bitget/merkel_tree_bg.json"))
        self.assertEqual(proc.returncode, 1)
        doc = json.loads(proc.stdout)
        self.assertEqual(doc["verdict"], "FAIL")
        self.assertFalse(doc["binding"]["matches"])
        self.assertIn("ROOT_NOT_BOUND", [f["id"] for f in doc["findings"]])

    def test_root_match_with_assessment_passes(self):
        root = json.loads((FIXTURES / "bitget/merkel_tree_bg.json").read_text())["root"]["hash"]
        assessment = self.write("bitget-ok.json", {
            "exchange": "bitget",
            "snapshot": {"id": "202605", "exchangeMerkleRoot": root},
        })
        proc = run(INCLUSION, "--exchange", "bitget", "--assessment", str(assessment),
                   str(FIXTURES / "bitget/merkel_tree_bg.json"))
        self.assertEqual(proc.returncode, 0, proc.stderr)
        self.assertTrue(json.loads(proc.stdout)["binding"]["matches"])

    def test_assessment_without_a_root_does_not_claim_a_match(self):
        assessment = self.write("no-root-assessment.json", {
            "exchange": "bitget",
            "snapshot": {"id": "202606", "exchangeMerkleRoot": ""},
        })
        proc = run(INCLUSION, "--exchange", "bitget", "--assessment", str(assessment),
                   str(FIXTURES / "bitget/merkel_tree_bg.json"))
        self.assertEqual(proc.returncode, 0, proc.stderr)
        binding = json.loads(proc.stdout)["binding"]
        self.assertIsNone(binding["matches"])
        self.assertIn("publishes no exchange merkle root", binding["note"])


class BindingToTheAnchor(unittest.TestCase):
    """Binding against a recorded on-chain read, with an honest live=false."""

    def test_anchor_root_binds_and_says_it_was_a_replay(self):
        proc = run(INCLUSION, "--exchange", "bybit",
                   "--onchain", "2025061709",
                   "--anchor-fixture", str(ANCHOR_FIXTURE),
                   str(FIXTURES / "bybit/mock_user_merkle_tree_path_40_v5.json"))
        self.assertEqual(proc.returncode, 0, proc.stderr)
        doc = json.loads(proc.stdout)
        self.assertTrue(doc["binding"]["matches"])
        self.assertFalse(doc["binding"]["anchor"]["live"])
        self.assertIn("anchor", doc["binding"]["publishedRoots"])

    def test_bitget_anchor_holds_a_64_bit_root_in_a_256_bit_word(self):
        """The contract stores the 64-bit root right-aligned; the scheme also emits 64 bits."""
        proc = run(INCLUSION, "--exchange", "bitget",
                   "--onchain", "202605",
                   "--anchor-fixture", str(ANCHOR_FIXTURE),
                   str(FIXTURES / "bitget/merkel_tree_bg.json"))
        doc = json.loads(proc.stdout)
        self.assertEqual(doc["verdict"], "PASS")
        self.assertIsNone(doc["binding"]["matches"])
        self.assertTrue(doc["binding"]["widthsDiffer"])
        self.assertIn("ROOT_WIDTH_MISMATCH", [f["id"] for f in doc["findings"]])

    def test_unavailable_anchor_is_not_a_verdict_about_the_proof(self):
        empty = json.loads(Path(ANCHOR_FIXTURE).read_text())
        empty["records"] = [row for row in empty["records"] if row["auditId"] != "2025061709"]
        tmp = Path(tempfile.mkdtemp(prefix="por-anchor-"))
        try:
            path = tmp / "no-record.json"
            path.write_text(json.dumps(empty))
            proc = run(INCLUSION, "--exchange", "bybit", "--onchain", "2025061709",
                       "--anchor-fixture", str(path),
                       str(FIXTURES / "bybit/mock_user_merkle_tree_path_40_v5.json"))
            doc = json.loads(proc.stdout)
            self.assertEqual(doc["verdict"], "PASS")
            self.assertIsNone(doc["binding"]["matches"])
            self.assertEqual(doc["binding"]["anchor"]["status"], "not_found")
        finally:
            shutil.rmtree(tmp, ignore_errors=True)


class NoSecretsEchoed(unittest.TestCase):
    def test_output_does_not_repeat_account_identifiers(self):
        proof = json.loads((FIXTURES / "mexc/merkle_proof.json").read_text())
        account = proof["self"]["accountId"]
        proc = run(INCLUSION, "--exchange", "mexc", str(FIXTURES / "mexc/merkle_proof.json"))
        self.assertNotIn(account, proc.stdout)
        self.assertNotIn(account, proc.stderr)

    def test_no_balances_appear_in_the_verdict(self):
        proof = json.loads((FIXTURES / "toobit/user_1_demo_proof.json").read_text())
        amounts = list(proof["self"]["balances"].values())
        proc = run(INCLUSION, "--exchange", "toobit", str(FIXTURES / "toobit/user_1_demo_proof.json"))
        for amount in amounts:
            self.assertNotIn(amount, proc.stdout)

    def test_nonce_is_never_echoed(self):
        proof = json.loads((FIXTURES / "coinex/mock_user_proof.json").read_text())
        proc = run(INCLUSION, "--exchange", "coinex", str(FIXTURES / "coinex/mock_user_proof.json"))
        self.assertNotIn(proof["self"]["nonce"], proc.stdout)

    def test_the_script_writes_no_output_file(self):
        before = sorted(p.name for p in (FIXTURES / "toobit").iterdir())
        run(INCLUSION, "--exchange", "toobit", str(FIXTURES / "toobit/user_1_demo_proof.json"))
        after = sorted(p.name for p in (FIXTURES / "toobit").iterdir())
        self.assertEqual(before, after)


class DecisionTable(unittest.TestCase):
    """Scores in decision-table.md must come out of the script unchanged."""

    # file, notes, expected score, expected level, anchor mode
    #
    # "anchor" rows use the recorded Sepolia reads in fixtures/anchor, so the row also
    # exercises the on-chain comparison. Rows where the real fixture contradicts the
    # file are scored file-only and say why.
    TABLE = [
        ("docs/reports/okx/508756245-assessment.json",
         "NO_CANONICAL_ANCHOR +5, HIGH_FREQUENCY_GAP +5, NO_PUBLIC_POR_API +5; not anchored, so file-only",
         15, "low", "--no-anchor"),
        ("docs/reports/kucoin/20261031-assessment.json",
         "UNVERIFIABLE +40, stage +35, list +16, ownership +8, NO_PUBLIC_POR_API (cap); anchored",
         100, "critical", "anchor"),
        ("docs/reports/bitget/20261021-assessment.json",
         "UNVERIFIABLE +40, stage +35, list +16, ownership +8, NO_PUBLIC_POR_API (cap); anchored",
         100, "critical", "anchor"),
        ("docs/reports/bybit/2026102706-assessment.json",
         "UNVERIFIABLE +40, stage +35, list +16, ownership +8, NO_PUBLIC_POR_API (cap)",
         100, "critical", "--no-anchor"),
        ("docs/reports/binance/PR04NOV26-assessment.json",
         "UNVERIFIABLE +40, stage +35, ownership +8, OPAQUE +8; has walletZip",
         91, "critical", "--no-anchor"),
        ("docs/reports/binance/PR01JUN26-assessment.json",
         "UNVERIFIABLE +40, stage +35, NO_WALLET_OWNERSHIP_PROOF +8, OPAQUE_TRUSTED_SETUP +8",
         91, "critical", "--no-anchor"),
        ("docs/reports/mexc/202610-assessment.json",
         "UNVERIFIABLE +40, stage +35, derived list +16, non-blocking capped; anchored",
         100, "critical", "anchor"),
    ]

    def test_no_wallet_list_scores_at_least_list_without_ownership(self):
        import risk_score

        bybit = risk_score.score(ROOT / "docs/reports/bybit/2026102706-assessment.json", None)
        binance = risk_score.score(ROOT / "docs/reports/binance/PR04NOV26-assessment.json", None)
        self.assertGreaterEqual(bybit["riskScore"], binance["riskScore"])

    def test_table_scores(self):
        import risk_score

        for path, _flags, expected_score, expected_level, mode in self.TABLE:
            with self.subTest(path=path):
                anchor = None if mode == "--no-anchor" else self._anchor(path)
                doc = risk_score.score(ROOT / path, anchor)
                self.assertTrue(doc["conclusive"])
                self.assertEqual(doc["riskScore"], expected_score)
                self.assertEqual(doc["riskLevel"], expected_level)
                self.assertNotIn("decision", doc)

    def _anchor(self, path):
        import por_anchor

        assessment = json.loads((ROOT / path).read_text())
        return por_anchor.resolve_record(
            str(assessment["snapshot"]["id"]),
            exchange=assessment["exchange"],
            fixture=str(ANCHOR_FIXTURE),
        )

    def test_blocking_flag_floors_the_level(self):
        doc = self._score({
            "exchange": "x", "porStage": "Stage 3", "effectivePoR": True,
            "riskFlags": [{"id": "LOGIN_WALL", "severity": "high"}],
        })
        self.assertEqual(doc["riskScore"], 40)
        self.assertEqual(doc["riskLevel"], "high")

    def test_stage_below_effective_por_floors_the_level(self):
        doc = self._score({
            "exchange": "x", "porStage": "Pre-Stage", "effectivePoR": False, "riskFlags": [],
        })
        self.assertEqual(doc["riskScore"], 35)
        self.assertEqual(doc["riskLevel"], "medium")

    def test_non_blocking_flags_are_capped_at_20(self):
        doc = self._score({
            "exchange": "x", "porStage": "Stage 2", "effectivePoR": True,
            "riskFlags": [{"id": "STAGE_BLOCKED"}, {"id": "NO_WALLET_OWNERSHIP_PROOF"},
                          {"id": "OPAQUE_TRUSTED_SETUP"}],
        })
        self.assertEqual(doc["riskScore"], 20)
        self.assertEqual([f["delta"] for f in doc["factors"]], [8, 8, 4])

    def test_unknown_flag_id_falls_back_to_severity(self):
        doc = self._score({
            "exchange": "x", "porStage": "Stage 2", "effectivePoR": True,
            "riskFlags": [{"id": "SOMETHING_NEW", "severity": "critical"}],
        })
        self.assertEqual(doc["riskScore"], 10)

    def test_unlisted_severity_defaults_to_5(self):
        doc = self._score({
            "exchange": "x", "porStage": "Stage 2", "effectivePoR": True,
            "riskFlags": [{"id": "SOMETHING_NEW", "severity": "weird"}],
        })
        self.assertEqual(doc["riskScore"], 5)

    def test_score_never_exceeds_100(self):
        doc = self._score({
            "exchange": "x", "porStage": "Stage 0", "effectivePoR": False,
            "riskFlags": [{"id": "UNVERIFIABLE"}, {"id": "UNVERIFIABLE"},
                          {"id": "STAGE_BLOCKED"}, {"id": "OPAQUE_TRUSTED_SETUP"}],
        })
        self.assertLessEqual(doc["riskScore"], 100)
        self.assertEqual(doc["riskLevel"], "critical")

    def test_second_run_lands_on_the_same_integers(self):
        import risk_score

        for path, _f, _s, _l, mode in self.TABLE:
            with self.subTest(path=path):
                anchor = None if mode == "--no-anchor" else self._anchor(path)
                first = risk_score.score(ROOT / path, anchor)
                second = risk_score.score(ROOT / path, anchor)
                self.assertEqual(json.dumps(first, sort_keys=True), json.dumps(second, sort_keys=True))

    def _score(self, assessment):
        import risk_score

        tmp = Path(tempfile.mkdtemp(prefix="por-score-"))
        try:
            path = tmp / "assessment.json"
            path.write_text(json.dumps(assessment))
            return risk_score.score(path, None)
        finally:
            shutil.rmtree(tmp, ignore_errors=True)


class AnchorShortCircuits(unittest.TestCase):
    """The on-chain read replaces the NEAR get; each outcome keeps its own meaning."""

    def setUp(self):
        sys.path.insert(0, str(ROOT / "skills" / "por-custody-risk" / "scripts"))
        import risk_score
        import por_anchor

        self.risk_score = risk_score
        self.por_anchor = por_anchor
        self.tmp = Path(tempfile.mkdtemp(prefix="por-anchor-cases-"))

    def tearDown(self):
        shutil.rmtree(self.tmp, ignore_errors=True)

    def _path(self, doc):
        path = self.tmp / "a.json"
        path.write_text(json.dumps(doc))
        return path

    def test_missing_file_is_100_critical(self):
        doc = self.risk_score.score(self.tmp / "absent.json", None)
        self.assertEqual(doc["riskScore"], 100)
        self.assertEqual(doc["riskLevel"], "critical")
        self.assertEqual(doc["factors"][0]["id"], "FILE_MISSING")

    def test_unreadable_anchor_is_inconclusive_not_low_risk(self):
        anchor = {"status": "unavailable", "live": False, "error": "all endpoints timed out"}
        doc = self.risk_score.score(self._path({
            "exchange": "x", "porStage": "Stage 2", "effectivePoR": True, "riskFlags": [],
        }), anchor)
        self.assertFalse(doc["conclusive"])
        self.assertIsNone(doc["riskScore"])
        self.assertIsNone(doc["riskLevel"])
        self.assertEqual(doc["factors"][0]["id"], "ANCHOR_UNAVAILABLE")

    def test_not_anchored_says_so_and_is_not_a_reserve_judgement(self):
        anchor = {"status": "not_found", "live": True, "note": "SnapshotNotFound"}
        doc = self.risk_score.score(self._path({
            "exchange": "x", "porStage": "Stage 2", "effectivePoR": True, "riskFlags": [],
        }), anchor)
        self.assertEqual(doc["factors"][0]["id"], "ANCHOR_MISSING")
        self.assertTrue(doc.get("notAReserveJudgement"))
        self.assertIn("not 'reserves are bad'", doc["factors"][0]["note"])

    def test_anchor_mismatch_reports_the_disagreeing_field(self):
        anchor = {
            "status": "found", "live": True,
            "record": {
                "artifactBundleRoot": "0x" + "11" * 32,
                "verificationBundleRoot": "0x" + "22" * 32,
                "exchangeMerkleRoot": "0x" + "33" * 32,
            },
        }
        doc = self.risk_score.score(self._path({
            "exchange": "x", "porStage": "Stage 2", "effectivePoR": True, "riskFlags": [],
            "artifacts": {"artifactBundleRoot": "0x" + "11" * 32,
                          "verificationBundleRoot": "0x" + "ff" * 32},
            "snapshot": {"id": "s1", "exchangeMerkleRoot": "33" * 32},
        }), anchor)
        self.assertEqual(doc["factors"][0]["id"], "ANCHOR_MISMATCH")
        self.assertIn("verificationBundleRoot", doc["factors"][0]["note"])

    def test_anchored_and_local_file_absent_both_sides_are_named(self):
        anchor = {
            "status": "found", "live": True,
            "record": {"artifactBundleRoot": "0x" + "11" * 32,
                       "verificationBundleRoot": "0x" + "22" * 32,
                       "exchangeMerkleRoot": "0x" + "00" * 32},
        }
        doc = self.risk_score.score(self.tmp / "absent.json", anchor)
        self.assertEqual(doc["factors"][0]["id"], "FILE_MISSING")
        self.assertIn("does hold an anchored record", doc["factors"][0]["note"])

    def test_zero_root_on_chain_is_not_compared(self):
        anchor = {
            "status": "found", "live": True,
            "record": {"artifactBundleRoot": "0x" + "11" * 32,
                       "verificationBundleRoot": "0x" + "22" * 32,
                       "exchangeMerkleRoot": "0x" + "00" * 32},
        }
        doc = self.risk_score.score(self._path({
            "exchange": "x", "porStage": "Stage 2", "effectivePoR": True, "riskFlags": [],
            "artifacts": {"artifactBundleRoot": "0x" + "11" * 32,
                          "verificationBundleRoot": "0x" + "22" * 32},
            "snapshot": {"id": "s1", "exchangeMerkleRoot": ""},
        }), anchor)
        rows = {row["field"]: row for row in doc["anchor"]["roots"]}
        self.assertFalse(rows["exchangeMerkleRoot"]["compared"])
        self.assertIsNone(rows["exchangeMerkleRoot"]["matches"])
        self.assertTrue(doc["conclusive"])

    def test_skipped_anchor_still_scores_the_file(self):
        doc = self.risk_score.score(ROOT / "docs/reports/okx/508756245-assessment.json",
                                    {"status": "skipped", "live": False})
        self.assertEqual(doc["riskScore"], 15)
        self.assertIn("file's own disclosure", doc["anchor"]["interpretation"])


class AnchorReader(unittest.TestCase):
    """The reader's Keccak and ABI decode, checked against the documented values."""

    def setUp(self):
        import por_anchor

        self.pa = por_anchor

    def test_keccak_matches_the_documented_exchange_tags(self):
        self.assertEqual(self.pa.exchange_tag("binance"),
                         "0x099e8cf4d817a6e4eec62bff4cdef05faa4a00fcde8d7e99f5090708d23ad9b2")
        self.assertEqual(self.pa.exchange_tag("okx"),
                         "0xd6b6d5e0aacce0469a313983d889ed10d0bb7c9545af0285a19b4ff094b4041d")

    def test_snapshot_ids_match_the_documented_values(self):
        self.assertEqual(self.pa.legacy_snapshot_id("PR01JUN26"),
                         "0x8eebf94be87651156d5d4b7eed3d99daed7fe308ba9ab0771cd5f6efb417c638")
        self.assertEqual(self.pa.legacy_snapshot_id("506872725"),
                         "0x88875070b19fd65cafbdd4674a55d835e0f24a96091375b001aef45e4248b2f0")

    def test_composite_snapshot_ids_avoid_audit_id_collisions(self):
        audit = "20261101"
        phemex = self.pa.composite_snapshot_id(audit, "phemex")
        htx = self.pa.composite_snapshot_id(audit, "htx")
        ourbit = self.pa.composite_snapshot_id(audit, "ourbit")
        legacy = self.pa.legacy_snapshot_id(audit)
        self.assertEqual(len({phemex, htx, ourbit, legacy}), 4)
        self.assertEqual(
            phemex,
            "0x12abd55db2e6d9b44edf24364ab8735226d612d54b228074210cb4933ec39f80",
        )

    def test_keccak_empty_string(self):
        self.assertEqual(self.pa.keccak_hex(b""),
                         "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")

    def test_record_decodes_from_a_recorded_return(self):
        doc = json.loads(ANCHOR_FIXTURE.read_text())
        row = next(r for r in doc["records"] if r["auditId"] == "PR01JUN26")
        record = row["record"]
        self.assertEqual(record["exchange"], "binance")
        self.assertEqual(record["periodSeq"], 43)
        self.assertEqual(record["schemaVersion"], 3)
        self.assertEqual(record["snapshotId"], self.pa.legacy_snapshot_id("PR01JUN26"))
        self.assertEqual(record["exchangeTag"], self.pa.exchange_tag("binance"))

    def test_empty_record_decodes_to_none(self):
        zeros = "0x" + "00" * 32 * 13
        self.assertIsNone(self.pa.decode_record(zeros))

    def test_snapshot_not_found_revert_is_recognised(self):
        selector = self.pa.selector("SnapshotNotFound(bytes32)")
        payload = selector + "11" * 32
        decoded = self.pa.decode_error(payload)
        self.assertEqual(decoded["name"], "SnapshotNotFound")

    def test_fixture_never_claims_to_be_live(self):
        doc = self.pa.resolve_record("PR01JUN26", exchange="binance",
                                     fixture=str(ANCHOR_FIXTURE))
        self.assertFalse(doc["live"])
        self.assertEqual(doc["status"], "found")

    def test_fixture_miss_is_reported_as_not_found(self):
        doc = self.pa.resolve_record("NOT-AN-ID", exchange="binance",
                                     fixture=str(ANCHOR_FIXTURE))
        self.assertEqual(doc["status"], "not_found")
        self.assertFalse(doc["live"])


class SkillDocExamples(unittest.TestCase):
    """The example run quoted in each SKILL.md must be what the script prints.

    The docs are the thing a reader copies, so a stale example is a bug in the
    skill. Each SKILL.md is parsed for its `$ ` command and the JSON block under
    it; the command is re-run and its output compared field by field against the
    quoted keys (the quotes are abridged where they say so).
    """

    SKILLS = {
        "por-custody-risk": CUSTODY_DIR / "SKILL.md",
        "por-inclusion-proof": INCLUSION_DIR / "SKILL.md",
    }

    def _quoted_example(self, path):
        """(command, quoted JSON object) from the first fenced block with a `$ ` line.

        The quoted run and its output share one block: the `$ ` lines are the
        command, a trailing backslash continues it, and the JSON after them is what
        the script printed. Fences are walked line by line so a closing fence is
        never mistaken for an opening one.
        """
        blocks = []
        current = None
        for line in path.read_text().splitlines():
            if line.startswith("```"):
                if current is None:
                    current = []
                else:
                    blocks.append(current)
                    current = None
                continue
            if current is not None:
                current.append(line)
        for lines in blocks:
            if not any(line.startswith("$ ") for line in lines):
                continue
            command_lines = []
            rest = []
            index = 0
            while index < len(lines):
                line = lines[index]
                if line.startswith("$ "):
                    command_lines.append(line[2:].strip())
                    while command_lines[-1].endswith("\\") and index + 1 < len(lines):
                        command_lines[-1] = command_lines[-1][:-1].strip()
                        index += 1
                        command_lines.append(lines[index].strip())
                    rest = lines[index + 1:]
                index += 1
            command = " ".join(part for part in command_lines if part)
            payload = "\n".join(rest).strip()
            if command and payload.startswith("{"):
                return command, json.loads(payload)
        self.fail(f"no quoted example run found in {path}")

    def test_custody_example_matches_the_script(self):
        command, quoted = self._quoted_example(self.SKILLS["por-custody-risk"])
        self.assertIn("por-custody-risk", command)
        # The quoted run uses the recorded fixture, so it is reproducible offline.
        proc = subprocess.run(command, shell=True, capture_output=True, text=True, cwd=ROOT)
        self.assertEqual(proc.returncode, 0, proc.stderr)
        actual = json.loads(proc.stdout)
        for key, expected in quoted.items():
            with self.subTest(key=key):
                self.assertIn(key, actual)
                if key == "anchor":
                    for anchor_key, anchor_value in expected.items():
                        self.assertEqual(actual["anchor"].get(anchor_key), anchor_value,
                                         f"anchor.{anchor_key} differs from the SKILL.md quote")
                else:
                    self.assertEqual(actual[key], expected, f"{key} differs from the SKILL.md quote")
        # The quote is a replay, and it must say so rather than look like a live read.
        self.assertFalse(actual["anchor"]["live"])

    def test_inclusion_example_matches_the_script(self):
        command, quoted = self._quoted_example(self.SKILLS["por-inclusion-proof"])
        self.assertIn("por-inclusion-proof", command)
        proc = subprocess.run(command, shell=True, capture_output=True, text=True, cwd=ROOT)
        self.assertEqual(proc.returncode, 0, proc.stderr)
        actual = json.loads(proc.stdout)
        for key, expected in quoted.items():
            with self.subTest(key=key):
                self.assertEqual(actual.get(key), expected, f"{key} differs from the SKILL.md quote")

    def test_listed_schemes_block_matches_the_script(self):
        text = (INCLUSION_DIR / "SKILL.md").read_text()
        block = re.search(r"\$ .*--list-schemes\n(.*?)```", text, re.S)
        self.assertIsNotNone(block, "SKILL.md no longer quotes --list-schemes")
        quoted = [line.strip() for line in block.group(1).splitlines() if line.strip()]
        actual = run(INCLUSION, "--list-schemes").stdout.split()
        self.assertEqual(quoted, actual)


class NoGoBuildRequired(unittest.TestCase):
    """These skills are meant to run on a bare python3 with no web3 library."""

    def test_scripts_import_only_the_standard_library(self):
        for script in (INCLUSION, CUSTODY):
            source = script.read_text()
            for banned in ("import web3", "from web3", "import requests", "import eth_",
                           "from eth_"):
                self.assertNotIn(banned, source, f"{script.name} imports {banned}")

    def test_lib_module_parses_under_a_plain_interpreter(self):
        proc = subprocess.run([sys.executable, "-c",
                               "import sys; sys.path.insert(0, 'skills/lib'); import por_anchor"],
                              capture_output=True, text=True, cwd=ROOT)
        self.assertEqual(proc.returncode, 0, proc.stderr)


if __name__ == "__main__":
    unittest.main()
