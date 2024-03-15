import asyncio

from grpclib.client import Channel

from ozonmp.logistic_package_api.v1.logistic_package_api_grpc import OmpTemplateApiServiceStub
from ozonmp.logistic_package_api.v1.logistic_package_api_pb2 import DescribeTemplateV1Request

async def main():
    async with Channel('127.0.0.1', 8082) as channel:
        client = OmpTemplateApiServiceStub(channel)

        req = DescribeTemplateV1Request(template_id=1)
        reply = await client.DescribeTemplateV1(req)
        print(reply.message)


if __name__ == '__main__':
    asyncio.run(main())
