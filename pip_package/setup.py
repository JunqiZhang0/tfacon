import os
import re

from setuptools import find_packages, setup

ROOT = os.path.dirname(__file__)
VERSION_RE = re.compile(r'''__version__ = ['"]([a-zA-Z0-9.]+)['"]''')


def get_version():
    init = open(os.path.join(ROOT, 'tfacon_pip', '__init__.py')).read()
    return VERSION_RE.search(init).group(1)


setup(
    name='tfacon',
    version=get_version(),
    description="tfacon",
    author="Red Hat Inc",
)
os.system('source scripts/install_tfacon.sh')