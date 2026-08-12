// === Drag Diagnostic Script ===
// Paste this in browser console (F12 -> Console) to check if drag is working

(function() {
    console.log('=== 🔍 Drag Functionality Check ===\n');
    
    // Check if vertical-splitter exists
    const splitter = document.querySelector('.vertical-splitter');
    if (!splitter) {
        console.error('❌ .vertical-splitter not found!');
        return;
    }
    
    console.log('✅ Found .vertical-splitter:', splitter);
    
    // Check computed styles
    const style = window.getComputedStyle(splitter);
    console.log('\n📊 Computed styles:');
    console.log('- width:', style.width);
    console.log('- display:', style.display);
    console.log('- position:', style.position);
    console.log('- z-index:', style.zIndex);
    console.log('- cursor:', style.cursor);
    console.log('- background-color:', style.backgroundColor);
    
    // Check if it's inside split-view-layout
    let parent = splitter.parentElement;
    let path = [];
    while (parent) {
        if (parent.classList && parent.classList.contains('split-view-layout')) {
            path.push('INSIDE split-view-layout ✅');
            break;
        }
        if (parent.tagName === 'EL-DIALOG') {
            path.push('INSIDE el-dialog ⚠️');
            break;
        }
        path.push(parent.tagName || parent.className);
        parent = parent.parentElement;
    }
    console.log('\n📂 Parent chain:', path.join(' > '));
    
    // Try to get event listeners (might not work in some browsers)
    try {
        // This won't show private listeners, but can check attached handlers
        console.log('\n⚠️ Note: Cannot check Vue event handlers directly from console');
        console.log('   Please check browser devtools Sources tab for the actual handler');
    } catch(e) {}
    
    // Test clickability
    const rect = splitter.getBoundingClientRect();
    console.log('\n📐 Position info:');
    console.log('- x:', Math.round(rect.x));
    console.log('- y:', Math.round(rect.y));
    console.log('- width:', Math.round(rect.width));
    console.log('- height:', Math.round(rect.height));
    
    // Visual test - temporarily highlight the element
    splitter.style.outline = '3px solid red';
    console.log('\n🎯 Element highlighted with RED border for visual confirmation');
    
})();
