export const prerender = true;

export async function GET() {
  console.log('SITEMAP GET FUNCTION CALLED!');
  const pages = [
    '',
    '/panduan',
    '/quiz',
    '/compress-image',
    '/berhitung',
    '/mengetik',
    '/cetak-kode',
  ];

  const website = 'https://lesbalonggarut.my.id';

  const sitemap = `<?xml version="1.0" encoding="UTF-8" ?>
<urlset
  xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xhtml="https://www.w3.org/1999/xhtml"
  xmlns:mobile="https://www.google.com/schemas/sitemap-mobile/1.0"
  xmlns:news="https://www.google.com/schemas/sitemap-news/0.9"
  xmlns:image="https://www.google.com/schemas/sitemap-image/1.1"
  xmlns:video="https://www.google.com/schemas/sitemap-video/1.1"
>
  ${pages
      .map(
        (page) => `  <url>
    <loc>${website}${page}</loc>
    <changefreq>daily</changefreq>
    <priority>${page === '' ? '1.0' : '0.7'}</priority>
  </url>`
      )
      .join('\n')}
</urlset>`;

  return new Response(sitemap, {
    headers: {
      'Cache-Control': 'max-age=0, s-maxage=3600',
      'Content-Type': 'application/xml',
    }
  });
}
