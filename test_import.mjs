import { readFileSync } from 'fs';

async function test() {
  try {
    console.log('Testing import of datastar.js');
    // Use file:// URL
    const module = await import('file:///root/reef/static/datastar.js?v=datastar-fix');
    console.log('Import succeeded');
    console.log('Exports:', Object.keys(module));
    console.log('mergePaths exists:', typeof module.mergePaths);
    console.log('mergePatch exists:', typeof module.mergePatch);
  } catch (err) {
    console.error('Import failed:', err);
    console.error('Stack:', err.stack);
  }
}

test();