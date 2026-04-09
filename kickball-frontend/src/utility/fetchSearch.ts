export async function fetchSearch(
   route: string,
   params?: URLSearchParams,
): Promise<any> {
   const query = params?.toString() ? `?${params.toString()}` : "";
   const url = `/${route}${query}`;
   const res = await fetch(url);
   if (!res.ok) throw new Error(`HTTP ${res.status} ${res.statusText}`);
   return res.json();
}
