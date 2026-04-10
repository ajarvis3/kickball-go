export async function fetchPost(
   route: string,
   headers: HeadersInit,
   body: BodyInit | null | undefined,
   params?: URLSearchParams,
) {
   const query = params?.toString() ? `?${params.toString()}` : "";
   // Read base URL from Vite env; default to empty string so behavior stays the same
   const rawBase = (import.meta.env.VITE_API_BASE ?? "") as string;
   const base = rawBase.replace(/\/+$/, ""); // strip trailing slashes
   // normalize route: remove leading slashes so prefix logic works consistently
   const cleanRoute = route.replace(/^\/+/, "");
   const prefix = base ? `${base}/` : "/";
   const url = `${prefix}${cleanRoute}${query}`;
   console.log(url);
   const res = await fetch(url, {
      method: "POST",
      headers: headers,
      body: body,
   });
   if (!res.ok) throw new Error(`HTTP ${res.status} ${res.statusText}`);
   return res.json();
}
