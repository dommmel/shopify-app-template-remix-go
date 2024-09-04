import { createChannel, createClient } from "nice-grpc";

import type { UserServiceClient } from "./generated/user";
import { UserServiceDefinition } from "./generated/user";

let _client: UserServiceClient;

function getClient(): UserServiceClient {
  if (!_client) {
    const address = process.env.API_URL;
    if (!address) {
      throw new Error("API URL is not set");
    }
    const channel = createChannel(address);

    _client = createClient(UserServiceDefinition, channel);
  }
  return _client;
}

export async function getUser(ID: number) {
  return await getClient().getUser({ID});
}

export async function findOrCreateUserByEncodedSessionToken(token: string) {
  return await getClient().findOrCreateUserByEncodedSessionToken({token});
}