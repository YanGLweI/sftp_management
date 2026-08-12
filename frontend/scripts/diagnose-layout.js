// === SFTP Browser Layout Diagnosis Script ===
// Run this in browser console (F12 -> Console)
// to check if split-view-layout is working correctly

(function diagnoseLayout() {
    console.log('=== 🚀 SFTP Browser Layout Diagnosis ===\n');
    
    // Check if dialog exists
    const dialog = document.querySelector('el-dialog');
    if (!dialog) {
        console.error('❌ No el-dialog found on page');
        return;
    }
    console.log('✅ el-dialog found:', dialog);
    
    // Check split-view-layout
    const layout = document.querySelector('.split-view-layout');
    if (!layout) {
        console.error('❌ No .split-view-layout found!');
        return;
    }
    console.log('✅ .split-view-layout found', layout);
    
    // Check computed styles
    const style = window.getComputedStyle(layout);
    console.log('\n📊 Computed styles:');
    console.log('- display:', style.display);
    console.log('- height:', style.height);
    console.log('- width:', style.width);
    console.log('- overflow:', style.overflow);
    console.log('- align-items:', style.alignItems);
    
    // Check flex children
    const children = Array.from(layout.children).map(el => ({
        tagName: el.tagName,
        className: el.className,
        width: window.getComputedStyle(el).width,
        display: window.getComputedStyle(el).display
    }));
    
    console.log('\n👥 Direct children:');
    children.forEach((child, idx) => {
        console.log(`[${idx}] <${child.tagName}>.${child.className.split(' ')[0]} | w:${child.width} | d:${child.display}`);
    });
    
    // Find panels specifically
    const leftPanel = document.querySelector('.left-panel');
    const rightPanel = document.querySelector('.right-panel');
    const splitter = document.querySelector('.vertical-splitter');
    
    console.log('\n🎯 Panels status:');
    console.log('- Left panel exists:', !!leftPanel);
    if (leftPanel) {
        const ls = window.getComputedStyle(leftPanel);
        console.log('  Width:', ls.width);
        console.log('  Display:', ls.display);
        console.log('  Position:', ls.position);
        console.log('  Flex:', ls.flex);
    }
    
    console.log('- Right panel exists:', !!rightPanel);
    if (rightPanel) {
        const rs = window.getComputedStyle(rightPanel);
        console.log('  Width:', rs.width);
        console.log('  Display:', rs.display);
        console.log('  Position:', rs.position);
        console.log('  Flex:', rs.flex);
        
        // Show rendered position
        const rect = rightPanel.getBoundingClientRect();
        console.log('  Rendered box: x=' + Math.round(rect.x) + ', y=' + Math.round(rect.y) + 
                    ', w=' + Math.round(rect.width) + ', h=' + Math.round(rect.height));
    } else {
        console.error('❌ RIGHT PANEL NOT FOUND - THIS IS THE PROBLEM!');
    }
    
    console.log('- Vertical splitter exists:', !!splitter);
    if (splitter) {
        const ss = window.getComputedStyle(splitter);
        console.log('  Width:', ss.width);
        console.log('  Display:', ss.display);
    }
    
    // Check if right panel is inside a dialog but outside split-view-layout
    const allRightPanels = document.querySelectorAll('.right-panel');
    console.log('\n📂 Total .right-panel elements:', allRightPanels.length);
    allRightPanels.forEach((rp, idx) => {
        let path = [];
        let parent = rp.parentElement;
        while (parent) {
            if (parent.classList && parent.classList.contains('split-view-layout')) {
                path.push('INSIDE split-view-layout');
                break;
            }
            if (parent.tagName === 'EL-DIALOG' || parent.id.includes('dialog')) {
                path.push('DIRECTLY IN DIALOG');
                break;
            }
            path.push(parent.tagName || 'div');
            parent = parent.parentElement;
        }
        console.log(`  [${idx}] Parent chain: ${path.join(' > ')}`);
    });
    
    // Check for common CSS issues
    console.log('\n⚠️  Common issues to check:');
    console.log('- If display is flex but items are vertical, check gap or margin');
    console.log('- If right panel is at bottom, it might not be a direct child of split-view-layout');
    console.log('- If panels are stacked horizontally instead of vertically, the layout container might be missing');
    
})();
console.log('💡 Execute: diagnoseLayout() to run diagnosis\n');
