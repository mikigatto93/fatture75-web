class CasingProductItem extends ProductItem {
    constructor(uuid, prodData, position) {
        super(uuid, prodData, position);
        this.depth = prodData.depth;
    }

    setupNode() {
        super.setupNode();
        
        this.node.querySelector(".casing-selector").style.display = "none";
    }
}