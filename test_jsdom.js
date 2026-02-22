const { JSDOM } = require('jsdom');
const fs = require('fs');
const path = require('path');

async function test() {
  // First, fetch the HTML from the local server
  const response = await fetch('http://localhost:3001/');
  const html = await response.text();
  
  // Create a virtual DOM with jsdom
  const dom = new JSDOM(html, {
    url: 'http://localhost:3001/',
    runScripts: 'dangerously',
    resources: 'usable',
    pretendToBeVisual: true,
    beforeParse(window) {
      // Mock fetch to capture debug beacons
      window.fetch = async (url, options) => {
        console.log('FETCH:', url);
        // Simulate success
        return {
          ok: true,
          json: async () => ({ status: 'ok' }),
          text: async () => '{}'
        };
      };
      
      // Capture console logs
      window.console.log = (...args) => console.log('LOG:', ...args);
      window.console.error = (...args) => console.error('ERROR:', ...args);
      
      // Mock EventSource
      window.EventSource = class MockEventSource {
        constructor(url) {
          console.log('EventSource created:', url);
          this.url = url;
          this.readyState = 0; // CONNECTING
          this.onopen = null;
          this.onerror = null;
          
          // Simulate connection opening after a delay
          setTimeout(() => {
            this.readyState = 1; // OPEN
            if (this.onopen) this.onopen();
          }, 10);
        }
        
        addEventListener() {}
        close() {}
      };
    }
  });
  
  const window = dom.window;
  
  // Wait a bit for scripts to execute
  await new Promise(resolve => setTimeout(resolve, 2000));
  
  // Check if Datastar loaded
  console.log('window.DatastarReady:', window.DatastarReady);
  console.log('window.datastarLoadFailed:', window.datastarLoadFailed);
  console.log('window.datastarLoadAttempts:', window.datastarLoadAttempts);
  console.log('window.Datastar:', window.Datastar ? 'exists' : 'missing');
  
  // Check for any errors
  window.addEventListener('error', (e) => {
    console.error('JS Error:', e.message);
  });
  
  // Wait a bit more
  await new Promise(resolve => setTimeout(resolve, 1000));
  
  console.log('Test complete');
}

test().catch(console.error);