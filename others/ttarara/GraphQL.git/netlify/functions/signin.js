export default async (req) => {
  if (req.method !== "POST") {
    return new Response("Method Not Allowed", { status: 405 });
  }

  const { username, password } = await req.json();

  if (!username || !password) {
    return new Response(JSON.stringify({ error: "Missing credentials" }), {
      status: 400,
      headers: { "Content-Type": "application/json" }
    });
  }

  const basic = Buffer.from(`${username}:${password}`).toString("base64");

  const res = await fetch("https://platform.zone01.gr/api/auth/signin", {
    method: "POST",
    headers: {
      Authorization: `Basic ${basic}`
    }
  });

  const token =
    res.headers.get("authorization")?.replace(/^Bearer\s+/i, "") ||
    res.headers.get("Authorization")?.replace(/^Bearer\s+/i, "");

  if (!res.ok || !token) {
    return new Response(JSON.stringify({ error: "Invalid credentials" }), {
      status: 401,
      headers: { "Content-Type": "application/json" }
    });
  }

  return new Response(JSON.stringify({ token }), {
    status: 200,
    headers: { "Content-Type": "application/json" }
  });
};
