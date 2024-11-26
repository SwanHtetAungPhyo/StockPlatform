import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '5s', target: 10 }, // Ramp-up to 10 virtual users in 5 seconds
    { duration: '10s', target: 10 }, // Hold at 10 virtual users for 10 seconds
    { duration: '5s', target: 0 },  // Ramp-down to 0 virtual users in 5 seconds
  ],
};

export default function () {
  const url = 'http://localhost/api/v1/'; // Change this URL to your API endpoint
  const res = http.get(url);

  // Check that the response status is 200
  check(res, {
    'is status 202': (r) => r.status === 202,
  });

  // Optionally: Sleep between requests to simulate real-world user interaction
  sleep(1); // Sleep for 1 second between requests
}
