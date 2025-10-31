import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { GreeterService } from "./api/v1/api_pb";

const client = createClient(
  GreeterService,
  createConnectTransport({
    baseUrl: window.location.origin,
  }),
);

document.getElementById("greet-btn")?.addEventListener("click", async () => {
  const response = await client.sayHello({ name: "World" });
  const responseDiv = document.getElementById("response");
  if (responseDiv) {
    responseDiv.textContent = response.message;
  }
});
