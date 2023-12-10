class CasingProductItem extends ProductItem {
    constructor(uuid, prodData) {
        super(uuid, prodData);
        this.depth = prodData.depth;
    }

    setupNode() {
        super.setupNode();
        
        this.node.querySelector(".casing-selector").style.display = "none";
    }

    
}