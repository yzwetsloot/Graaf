import json

import matplotlib.pyplot as plt
import networkx as nx
from networkx import drawing


with open('config/config.json') as config:
    parameters = json.load(config)

G = nx.DiGraph()


def main():
    with open('urls.txt', encoding='utf8') as fh:
        for line in fh:
            parse(line)

    nodelist = [node for node in G.nodes if G.in_degree(
        node) > parameters['min_indeg']]
    edgelist = [edge for edge in G.edges if edge[0]
                in nodelist and edge[1] in nodelist]

    labels = {node: node for node in nodelist}

    node_size = [G.in_degree(node) * parameters['node_width']
                 for node in nodelist]

    plot(nodelist, edgelist, labels, node_size)


def parse(line):
    components = line.strip().split(';')
    G.add_node(components[0])

    for opposite in neighbours(components[1]):
        G.add_edge(opposite, components[0])

    for opposite in neighbours(components[2]):
        G.add_edge(components[0], opposite)


def neighbours(line):
    if line is not '':
        return line.split(',')
    return []


def plot(nodelist=None, edgelist=None, labels=None, node_size=None):
    options = {
        'with_labels': True,
        'nodelist': nodelist,
        'edgelist': edgelist,
        'labels': labels,

        'font_size': parameters['font_size'],
        'font_color': parameters['font_color'],
        'font_weight': parameters['font_weight'],

        'node_size': node_size,

        'width': parameters['width'],
        'arrowsize': parameters['arrowsize'],

        'alpha': parameters['alpha'],
        'edge_color': parameters['edge_color'],
        'node_color': parameters['node_color']
    }

    nx.draw(G, nx.spring_layout(G), **options)

    plt.savefig('graph.svg')


if __name__ == "__main__":
    main()
