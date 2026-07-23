#!/usr/bin/env python3
import json
from pathlib import Path


BASE = Path(__file__).resolve().parent
GRAPH_PATH = BASE / "graph.json"
OVERRIDES_PATH = BASE / "manual_overrides.json"


def load_json(path):
    return json.loads(path.read_text(encoding="utf-8"))


def write_json(path, data):
    path.write_text(json.dumps(data, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def node_key(node):
    return (
        node.get("id"),
        node.get("label"),
        node.get("source_file"),
        node.get("source_location"),
    )


def find_node(nodes, spec):
    for node in nodes:
        if node.get("id") == spec["id"]:
            return node
    for node in nodes:
        if (
            node.get("label") == spec.get("label")
            and node.get("source_file") == spec.get("source_file")
            and node.get("source_location") == spec.get("source_location")
        ):
            return node
    return None


def make_node(spec):
    node = dict(spec)
    node.setdefault("_origin", "manual")
    return node


def link_key(link):
    return (link.get("source"), link.get("target"), link.get("relation"))


def main():
    graph = load_json(GRAPH_PATH)
    overrides = load_json(OVERRIDES_PATH)

    nodes = graph["nodes"]
    links = graph["links"]

    added_nodes = []
    for spec in overrides.get("add_nodes", []):
        if find_node(nodes, spec) is None:
            nodes.append(make_node(spec))
            added_nodes.append(spec["id"])

    existing = {link_key(link) for link in links}
    removed = 0
    remove_specs = overrides.get("remove_links", [])
    remove_keys = {(spec["source"], spec["target"], spec["relation"]) for spec in remove_specs}
    filtered = []
    for link in links:
        if link_key(link) in remove_keys:
            removed += 1
            continue
        filtered.append(link)
    links = filtered
    existing = {link_key(link) for link in links}

    added_links = []
    for spec in overrides.get("add_links", []):
        key = (spec["source"], spec["target"], spec["relation"])
        if key in existing:
            continue
        links.append(dict(spec))
        existing.add(key)
        added_links.append(key)

    graph["nodes"] = nodes
    graph["links"] = links
    graph.setdefault("hyperedges", [])

    write_json(GRAPH_PATH, graph)

    summary = {
        "added_nodes": added_nodes,
        "added_links": ["|".join(x) for x in added_links],
        "removed_links": removed,
        "node_count": len(nodes),
        "link_count": len(links),
    }
    write_json(BASE / "manual_overrides.applied.json", summary)
    print(json.dumps(summary, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
