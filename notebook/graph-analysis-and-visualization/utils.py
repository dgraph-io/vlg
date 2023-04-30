from IPython.display import Markdown, display

def warning(text):
    return Markdown(f'<div class="alert alert-danger">{text}</div>')
