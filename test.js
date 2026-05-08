import { randomItem } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { check } from 'k6';
import http from 'k6/http';

const BASE_URL = 'http://localhost:3006';
  const codesGlob = [];

export const options = {
      noConnectionReuse: false,
 setupTimeout: '60s',
  scenarios: {
    redirect: {
      executor: 'constant-arrival-rate',

      exec: 'visitShortUrl',

      rate: 1000,
      timeUnit: '1s',

      duration: '30s',

      preAllocatedVUs: 500,
      maxVUs: 2000,
    },
   load_test: {
      executor: 'ramping-arrival-rate',

      startRate: 1000,
      timeUnit: '1s',

      preAllocatedVUs: 500,
      maxVUs: 5000,

      stages: [
        { target: 1000, duration: '40s' },
          { target: 5000, duration: '20s' },
        { target: 100, duration: '10s' },
      ],
    },
}
};

export function setup() {
  const codes = [];

  for (let i = 0; i < 1000; i++) {
    const payload = JSON.stringify({
      url: 'https://example.com/',
    });

    const params = {
      headers: {
        'Content-Type': 'application/json',
      },
      tags: {
      name: 'POST /shorten',
    },
    };

    const res = http.post(
      `${BASE_URL}/shorten`,
      payload,
      params
    );

    check(res, {
        'GET status is 201': (r) => r.status === 201,
    });

    const shortUrl = res.body.trim();

    if (shortUrl) {
      codes.push(shortUrl);
    }
  }

  console.log(`Prepared ${codes.length} short urls`);

  return { codes };
}

export function visitShortUrl(data) {
  const code = randomItem(data.codes);

  const res = http.get(`${BASE_URL}/${code}`, {
  redirects: 0,
  tags: {
    name: 'GET /:code',
  },
});

  check(res, {
    'GET status is 301': (r) => r.status === 301,
  });
}

export function createShortUrl() {
  const payload = JSON.stringify({
    url: 'https://example.com/',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
    tags: {
      name: 'POST /shorten',
    },
  };

  const res = http.post(
    `${BASE_URL}/shorten`,
    payload,
    params,
  );

  if (!res || !res.body) {
    return;
    }

  check(res, {
    'POST status is 201': (r) => r.status === 201,
  });

  // например body:
  // http://localhost:3006/abc123

  const shortUrl = res.body.trim();

  const code = shortUrl.split('/').pop();

  if (code) {
    codesGlob.push(code);
  }
}

export default function () {
  const shouldPost = Math.random() < 0.5;

  if (shouldPost || codesGlob.length === 0) {
    createShortUrl();
  } else {
    visitShortUrl({codes:codesGlob});
  }
}

