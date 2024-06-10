import ctypes
from pathlib import Path

here = Path(__file__).absolute().parent
so_file = here / ("lego.so")
library = ctypes.cdll.LoadLibrary(so_file)
library.run.restype = ctypes.c_char_p


def hello():
    library.Hello()


def get_certificate(email: str, server: str, csr: bytes, plugin: str, env: dict[str, str]) -> str:
    envvar = list(sum(env.items(), ()))
    envs = (ctypes.c_char_p * len(envvar))()
    envs[:] = [bytes(e, "utf-8") for e in envvar]

    ctypes.c_char_p(bytes(email, "utf-8"))
    c_certificate_bundle: bytes = library.run(
        ctypes.c_char_p(bytes(email, "utf-8")),
        ctypes.c_char_p(bytes(server, "utf-8")),
        csr,
        ctypes.c_char_p(bytes(plugin, "utf-8")),
        envs,
        len(envs),
    )
    certificate_bundle = c_certificate_bundle.decode(encoding="utf-8")
    return certificate_bundle
