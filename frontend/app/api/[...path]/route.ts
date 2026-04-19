import { NextRequest, NextResponse } from "next/server";

const backendBaseURL = process.env.BACKEND_INTERNAL_URL || "http://api-gateway:8080";

type RouteParams = { path: string[] };
type RouteContext = { params: RouteParams | Promise<RouteParams> };

async function proxy(request: NextRequest, params: RouteParams) {
  const path = params.path.join("/");
  const search = request.nextUrl.search || "";
  const url = `${backendBaseURL}/api/${path}${search}`;

  const headers = new Headers(request.headers);
  headers.delete("host");

  const init: RequestInit = {
    method: request.method,
    headers,
    body: request.method === "GET" || request.method === "HEAD" ? undefined : await request.arrayBuffer(),
    cache: "no-store",
  };

  const response = await fetch(url, init);
  const body = await response.arrayBuffer();

  const responseHeaders = new Headers(response.headers);
  responseHeaders.delete("content-encoding");
  responseHeaders.delete("content-length");

  return new NextResponse(body, {
    status: response.status,
    headers: responseHeaders,
  });
}

export async function GET(request: NextRequest, context: RouteContext) {
  const params = await context.params;
  return proxy(request, params);
}

export async function POST(request: NextRequest, context: RouteContext) {
  const params = await context.params;
  return proxy(request, params);
}

export async function PATCH(request: NextRequest, context: RouteContext) {
  const params = await context.params;
  return proxy(request, params);
}

export async function DELETE(request: NextRequest, context: RouteContext) {
  const params = await context.params;
  return proxy(request, params);
}

export async function PUT(request: NextRequest, context: RouteContext) {
  const params = await context.params;
  return proxy(request, params);
}

export async function OPTIONS(request: NextRequest, context: RouteContext) {
  const params = await context.params;
  return proxy(request, params);
}

