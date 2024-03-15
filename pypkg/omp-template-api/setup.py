import setuptools

setuptools.setup(
    name="grpc-logistic-package-api",
    version="1.0.0",
    author="rusdevop",
    author_email="rusdevops@gmail.com",
    description="GRPC python client for omp-template-api",
    url="https://github.com/arslanovdi/logistic-package-api",
    packages=setuptools.find_packages(),
    package_data={"ozonmp.logistic_package_api.v1": ["logistic_package_api_pb2.pyi"]},
    python_requires='>=3.5',
)