class ProductItem {
    constructor(uuid, prodData) {
        this.uuid = uuid;
        this.prodData = prodData;
        this.node = null;
    }

    setupNode() {
        let templateNode = document.querySelector("#prod-item-template")
            .content
            .querySelector(".prod-list-item"); //get the li element

        this.node = templateNode.cloneNode(true);

        //setup event handlers
        let self = this;
        let checkbox = this.node.querySelector("#tapparelle-checkbox");
        checkbox.addEventListener(
            "change", function () {
                self.handleTapparelleCheckboxChange.call(self, checkbox.checked);
            }
        );
    }

    handleTapparelleCheckboxChange(checked) {
        if (checked) {
            console.log(this.uuid);
            this.node.querySelector("#prod-group").value = "B";
        }
    }

    render(parent) {
        this.node.querySelector(".width").textContent = this.prodData.width;
        this.node.querySelector(".height").textContent = this.prodData.height;
        this.node.querySelector(".prod-id").textContent = this.prodData.product_id;

        parent.appendChild(this.node);
    }

}