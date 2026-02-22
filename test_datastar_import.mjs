import('./static/datastar.js').then(module => {
  console.log('Datastar module loaded successfully');
  console.log('Exports:', Object.keys(module));
  console.log('mergePatch:', typeof module.mergePatch);
  console.log('mergePaths:', typeof module.mergePaths);
}).catch(err => {
  console.error('Failed to import datastar:', err);
  console.error(err.stack);
});