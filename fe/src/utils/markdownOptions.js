import React, { useEffect, useRef } from 'react';
import hljs from 'highlight.js';
import 'highlight.js/styles/github-dark.css'; 

const CodeBlock = ({ children }) => {
    const codeRef = useRef(null);

    useEffect(() => {
        if (codeRef.current) {
            hljs.highlightBlock(codeRef.current);
        }
    }, [children]);

    const codeProps = React.Children.only(children).props;
    const language = codeProps.className ? codeProps.className.replace('lang-', '') : 'plaintext';
    const codeText = codeProps.children;

    return (
        <pre className="p-4 rounded-md bg-gray-800 text-white">
            <code ref={codeRef} className={`language-${language}`}>
                {codeText}
            </code>
        </pre>
    );
};


export const baseMarkdownOptions = {
    overrides: {
        h1: {
            component: "h1",
            props: { style: { fontSize: "1.6rem", margin: "12px 0", fontWeight: "bold", textAlign: "center", color: "#f8e187" } },
        },
        h2: {
            component: "h2",
            props: { style: { fontSize: "1.4rem", margin: "10px 0", fontWeight: "bold", textAlign: "center", color: "#f8e187" } },
        },
        h3: {
            component: "h3",
            props: { style: { fontSize: "1.2rem", margin: "8px 0", fontWeight: "bold", textAlign: "center",  color: "#f8e187" } },
        },
        p: {
            component: "p",
            props: { style: { fontSize: "1rem", lineHeight: "1.6", margin: "10px 0", textAlign: "left" } },
        },
        ul: {
            component: "ul",
            props: { style: { paddingLeft: "20px", marginBottom: "12px", marginLeft: "10px", listStyleType: "circle", textAlign: "left" } },
        },
        li: {
            component: "li",
            props: { style: { marginBottom: "6px", lineHeight: "1.5", marginLeft: "15px", textAlign: "left" } },
        },
        blockquote: {
            component: "blockquote",
            props: { style: { fontSize: "1.1rem", margin: "10px 0", paddingLeft: "15px", borderLeft: "4px solid #ccc", fontStyle: "italic" } },
        },
        strong: {
            component: "strong",
            props: { style: { fontWeight: "bold" } },
        },
        em: {
            component: "em",
            props: { style: { fontStyle: "italic", color: "#555" } },
        },
        a: {
            component: "a",
            props: { style: { color: "#1e90ff", textDecoration: "none" } },
        },
        hr: {
            component: "hr",
            props: { style: { border: "0", height: "1px", backgroundColor: "#ccc", margin: "15px 0" } },
        },
        pre: {
            component: CodeBlock,
        },
    },
};