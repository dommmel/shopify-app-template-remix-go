import { HeadersFunction, MetaFunction, json, ActionFunctionArgs, LoaderFunctionArgs } from "@remix-run/node";
import { getUser } from "../grpc/user.server";
import { useRouteError, useLoaderData, useFetcher } from "@remix-run/react";
import { boundary } from "@shopify/shopify-app-remix/server";
import { authenticate } from "../.server/authenticate";
export const meta: MetaFunction = () => {
  return [
    { title: "New Remix App" },
    { name: "description", content: "Welcome to Remix!" },
  ];
};

export const loader = async ({request}:LoaderFunctionArgs) => {
  const authenticatedUser = await authenticate(request);

  return json(authenticatedUser);
};

export const action = async ({ request }:ActionFunctionArgs) => {
  await authenticate(request);
  const formData = await request.formData();
  const userId = formData.get("userId")?.toString() || "1";
  const user = await getUser(parseInt(userId, 10));
  return json(user);
};

export default function Index() {
  const loaderDataUser = useLoaderData<typeof loader>();
  const fetcher = useFetcher<typeof action>();

  // Check if has data, otherwise use initial data
  const user = fetcher.data || loaderDataUser;

  return (
    <div style={{ fontFamily: "system-ui, sans-serif", lineHeight: "1.8" }}>
      <h1>User</h1>
      <fetcher.Form method="post">
        <input type="number" name="userId" defaultValue={1} />
        <button type="submit">Fetch User</button>
      </fetcher.Form>
      <pre>{JSON.stringify(user, null, 2)}</pre>
      <pre>Current User: {JSON.stringify(loaderDataUser, null, 2)}</pre>
    </div>
  );
}

// Shopify needs Remix to catch some thrown responses, so that their headers are included in the response.
export function ErrorBoundary() {
  return boundary.error(useRouteError());
}

export const headers: HeadersFunction = (headersArgs) => {
  return boundary.headers(headersArgs);
};
