// js/utils/markdown.js 
export const processMarkdown = (text) => {
    const codeBlocks = [];
    const formulas = [];
    let placeholderIndex = 0;

    // === Шаг 1: Извлекаем ФОРМУЛЫ (сначала!) ===
    let processedText = text;

    // Блочные формулы: $$ ... $$
    processedText = processedText.replace(/\$\$(.*?)\$\$/gs, (match, formula) => {
        const index = formulas.length;
        formulas[index] = { type: 'block', content: formula.trim() };
        return `__FORMULA_${index}__`;
    });

    // Инлайн формулы: $ ... $
    processedText = processedText.replace(/\$(.*?)\$/g, (match, formula) => {
        const index = formulas.length;
        formulas[index] = { type: 'inline', content: formula.trim() };
        return `__FORMULA_${index}__`;
    });

    // === Шаг 2: Извлекаем КОДОВЫЕ БЛОКИ ===
    processedText = processedText.replace(/```(\w+)?\s*([\s\S]*?)\s*```/g, (match, lang, code) => {
        const idx = placeholderIndex;
        const langAttr = lang ? ` class="language-${lang}"` : ' class="language-plaintext"';
        // Важно: не экранируем содержимое кода!
        codeBlocks[idx] = `<pre class="no-math"><code${langAttr}>${code.trim()}</code></pre>`;
        placeholderIndex++;
        return `__CODE_BLOCK_${idx}__`;
    });

    // === Шаг 3: Экранируем ОСТАТОК текста ===
    processedText = escapeHtml(processedText);

    // === Шаг 4: Преобразуем Markdown в HTML ===
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
        .replace(/\n/g, '<br>'); // ← Переносы ПОСЛЕ обработки

    // === Шаг 5: Восстанавливаем кодовые блоки ===
    codeBlocks.forEach((block, i) => {
        processedText = processedText.replace(`__CODE_BLOCK_${i}__`, block);
    });

    // === Шаг 6: Восстанавливаем формулы ===
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

const escapeHtml = (str) => {
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '<')
        .replace(/>/g, '>')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
};