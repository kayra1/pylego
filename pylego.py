#!/usr/bin/env python3
import ctypes
from distutils.sysconfig import get_config_var
from pathlib import Path

here = Path(__file__).absolute().parent
ext_suffix = get_config_var('EXT_SUFFIX')
so_file = here / ('_pylego' + ext_suffix)
library = ctypes.cdll.LoadLibrary(so_file)
library.run.restype = ctypes.c_char_p


def run(email: str, server: str, csr: bytes, plugin: str, env: dict[str, str]) -> str:
    envvar = list(sum(env.items(), ()))
    envs = (ctypes.c_char_p * len(envvar))()
    envs[:] = [bytes(e, 'utf-8') for e in envvar]

    ctypes.c_char_p(bytes(email, 'utf-8'))
    c_certificate_bundle: bytes = library.run(
        ctypes.c_char_p(bytes(email, 'utf-8')),
        ctypes.c_char_p(bytes(server, 'utf-8')),
        csr,
        ctypes.c_char_p(bytes(plugin, 'utf-8')),
        envs,
        len(envs)
    )
    certificate_bundle = c_certificate_bundle.decode(encoding='utf-8')
    return certificate_bundle
