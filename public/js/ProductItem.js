class ProductItem {
    constructor(uuid, prodData) {
        this.uuid = uuid;
        this.position = prodData.position;
        this.prodData = prodData;
        this.node = null;
        this.group = "A";  // default group
    }

    //custom event system
    on(name, callback) {
        var callbacks = this[name];
        if (!callbacks) this[name] = [callback];
        else callbacks.push(callback);
    }

    dispatch(event) {
        var callbacks = this[event.name];
        if (callbacks) callbacks.forEach(callback => callback(event));
    }

    setupNode() {
        let templateNode = document.querySelector("#prod-item-template")
            .content
            .querySelector(".prod-list-item"); //get the li element

        this.node = templateNode.cloneNode(true);

        //setup event handlers
        let self = this;
        let groupChangeSelector = this.node.querySelector(".prod-group");
        
        groupChangeSelector.addEventListener(
            "change", function () {
                self.setGroup.call(self, groupChangeSelector.value);
            }
        );

    }

    setGroup(group){
        this.node.querySelector(".prod-group").value = group;
        this.group = group;
    }


    render(parent) {
        this.node.querySelector(".width").textContent = this.prodData.width;
        this.node.querySelector(".height").textContent = this.prodData.height;
        this.node.querySelector(".prod-id").textContent = this.prodData.product_id;

        this.node.querySelector(".position").textContent = this.position;
        this.node.querySelector(".quantity").textContent = this.prodData.quantity;
        
        this.node.querySelector(".reference").textContent = this.prodData.reference;

        this.node.querySelector(".tot-price").textContent = this.prodData.tot_price;

        parent.appendChild(this.node);
    }

    toJson() {
        return {
            "product_data": this.prodData,
            "position": parseInt(this.position),
            "group": this.group,
        }
    }

}