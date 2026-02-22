const fs = require('fs');
const vm = require('vm');

// Read dashboard.html
const html = fs.readFileSync('/root/reef/templates/dashboard.html', 'utf8');

// Extract all script tags content (non-module)
const scriptRegex = /<script(?:\s+[^>]*)?>(.*?)<\/script>/gs;
let match;
let scripts = [];
while ((match = scriptRegex.exec(html)) !== null) {
  // Skip module scripts and empty scripts
  if (!match[0].includes('type="module"') && match[1].trim()) {
    scripts.push(match[1]);
  }
}

console.log(`Found ${scripts.length} inline scripts`);

// Create a sandbox with mock browser APIs
const sandbox = {
  window: {},
  document: {
    head: {
      appendChild: () => {},
    },
    createElement: (tag) => ({
      src: '',
      crossOrigin: '',
      onload: null,
      onerror: null,
      type: '',
      textContent: '',
    }),
  },
  console: {
    log: (...args) => console.log('SCRIPT LOG:', ...args),
    error: (...args) => console.error('SCRIPT ERROR:', ...args),
    warn: (...args) => console.warn('SCRIPT WARN:', ...args),
  },
  fetch: (url) => {
    console.log('FETCH called:', url);
    return Promise.resolve({
      ok: true,
      json: () => Promise.resolve({}),
      text: () => Promise.resolve('{}'),
    });
  },
  location: {
    pathname: '/',
  },
  setTimeout,
  clearTimeout,
};

sandbox.window = sandbox;
sandbox.self = sandbox;

// Evaluate each script
for (let i = 0; i < scripts.length; i++) {
  const script = scripts[i];
  console.log(`\n=== Evaluating script ${i + 1} ===`);
  try {
    vm.createContext(sandbox);
    vm.runInContext(script, sandbox);
    console.log('Script executed without syntax error');
  } catch (err) {
    console.error('Syntax error in script:', err.message);
    console.error('Script snippet:', script.substring(0, 200));
  }
}

// Check if datastar loading variables were set
console.log('\n=== Sandbox state ===');
console.log('window.datastarLoadFailed:', sandbox.window.datastarLoadFailed);
console.log('window.datastarLoadAttempts:', sandbox.window.datastarLoadAttempts);
console.log('window.datastarLoaded:', sandbox.window.datastarLoaded);
console.log('window.DatastarReady:', sandbox.window.DatastarReady);