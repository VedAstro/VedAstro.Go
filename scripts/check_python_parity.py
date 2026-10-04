"""Compare generated Go metadata to Python without importing its side-effectful package."""
import ast
import json
import pathlib
import sys

root = pathlib.Path(__file__).resolve().parents[1]
source = pathlib.Path(sys.argv[1])
tree = ast.parse(source.read_text(encoding="utf-8-sig"))
calculate = next(node for node in tree.body if isinstance(node, ast.ClassDef) and node.name == "Calculate")
manifest = json.loads((root / "api_manifest.json").read_text(encoding="utf-8"))
python_methods = {node.name: node for node in calculate.body if isinstance(node, ast.FunctionDef)}
expected_names = set()
known_python_empty_string_bugs = {
    (f"Chapter{chapter}PrashnaMargaPredictions", "firstLetterOfQuery")
    for chapter in (13, 16, 25, 26)
}
observed_differences = set()
for method in manifest["methods"]:
    name = method["name"]
    expected_names.add(name)
    node = python_methods[name]
    args = node.args.args[1:]  # cls
    assert [arg.arg for arg in args] == [p["name"] for p in method["parameters"]], name
    defaults = {arg.arg: ast.literal_eval(default) for arg, default in zip(args[len(args)-len(node.args.defaults):], node.args.defaults)}
    for parameter in method["parameters"]:
        assert parameter["optional"] == (parameter["name"] in defaults), (name, parameter)
        if parameter["optional"]:
            actual, expected = defaults[parameter["name"]], parameter["defaultValue"]
            if actual != expected:
                key = (name, parameter["name"])
                assert key in known_python_empty_string_bugs and actual is None and expected == "", (name, parameter, defaults)
                observed_differences.add(key)
# Select calculations by their endpoint assignment, excluding configuration helpers.
actual_names = {name for name,node in python_methods.items() if any(
    isinstance(stmt, ast.Assign) and any(isinstance(target, ast.Name) and target.id == "endpoint" for target in stmt.targets)
    for stmt in node.body)}
assert actual_names == expected_names, (actual_names - expected_names, expected_names - actual_names)
print(f"Python parity passed: {len(expected_names)} methods, parameter order, optional flags, and defaults.")
if observed_differences:
    print(f"Documented exception: {len(observed_differences)} Python None defaults are empty strings in the server metadata; Go leaves them to the server.")
