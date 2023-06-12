# -*- coding: utf-8 -*-
# vim: set ft=python ts=4 sw=4 expandtab:

import trustero_api.receptor_v1.receptor_pb2 as receptor


class TestProtoc:
    def test_import_ok(self):
        assert receptor is not None
