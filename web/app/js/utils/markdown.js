// js/utils/markdown.js 

const escapeHtml = (str) => {
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '<')
        .replace(/>/g, '>')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
};

const processMarkdown = (text) => {
    const codeBlocks = [];
    const formulas = [];
    let placeholderIndex = 0;

    let processedText = text;

    processedText = processedText.replace(/\$\$(.*?)\$\$/gs, (match, formula) => {
        const index = formulas.length;
        formulas[index] = { type: 'block', content: formula.trim() };
        return `__FORMULA_${index}__`;
    });

    processedText = processedText.replace(/\$(.*?)\$/g, (match, formula) => {
        const index = formulas.length;
        formulas[index] = { type: 'inline', content: formula.trim() };
        return `__FORMULA_${index}__`;
    });

    processedText = processedText.replace(/```(\w+)?\s*([\s\S]*?)\s*```/g, (match, lang, code) => {
        const idx = placeholderIndex;
        const langAttr = lang ? ` class="language-${lang}"` : ' class="language-plaintext"';
        codeBlocks[idx] = `<pre class="no-math"><code${langAttr}>${code.trim()}</code></pre>`;
        placeholderIndex++;
        return `__CODE_BLOCK_${idx}__`;
    });

    processedText = escapeHtml(processedText);

    processedText = processedText
        .replace(/^###### (.*$)/gm, '<h6>$1</h6>')
        .replace(/^##### (.*$)/gm, '<h5>$1</h5>')
        .replace(/^#### (.*$)/gm, '<h4>$1</h4>')
        .replace(/^### (.*$)/gm, '<h3>$1</h3>')
        .replace(/^## (.*$)/gm, '<h2>$1</h2>')
        .replace(/^# (.*$)/gm, '<h1>$1</h1>')
        .replace(/^> (.*)(\n> .*)*/gm, (match, firstLine, rest = '') => {
            const lines = [firstLine, ...rest.split('\n> ').filter(l => l.trim())];
            return `<blockquote>${lines.join('<br>')}</blockquote>`;
        })
        .replace(/^---+$/gm, '<hr>')
        .replace(/`(.*?)`/g, '<code>$1</code>')
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
        .replace(/\n/g, '<br>');

    codeBlocks.forEach((block, i) => {
        processedText = processedText.replace(`__CODE_BLOCK_${i}__`, block);
    });

    formulas.forEach((formula, i) => {
        const placeholder = `__FORMULA_${i}__`;
        const { type, content } = formula;
        if (type === 'block') {
            processedText = processedText.replace(placeholder, `$$${content}$$`);
        } else {
            processedText = processedText.replace(placeholder, `$${content}$`);
        }
    });

    return processedText;
};

window.processMarkdown = processMarkdown;