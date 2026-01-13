import Markdown from "markdown-to-jsx";
import { baseMarkdownOptions } from "../utils/markdownOptions";

const MarkDown = ({ text = "" }) => {
    
    return (
        <Markdown
            options={baseMarkdownOptions}
        >
            {text}
        </Markdown>
    )
}

export default MarkDown;