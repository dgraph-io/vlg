## VLG Jupyter Notebook Resources

This folder houses various Jupyter Notebooks for working with the VLG. As we use [pipenv](https://pipenv.pypa.io/en/latest/) to bring some sanity to the python dev flow, each notebook is located in its own subfolder, with a corresponding Pipfile.

The only requirements for _all_ notebooks is Python 3.8 or later, and `pipenv`.

To run a specific notebook, `cd` to that folder, check out the README in that folder (each notebook has varying requirements), then:

```
pipenv install
pipenv run jupyter lab notebook.ipynb
```
