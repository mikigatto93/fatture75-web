class FixtureProductItem extends ProductItem {
    constructor(uuid, prodData) {
        super(uuid, prodData);
        this.rollerShutterPrice = 0;
        this.casing = null;
    }

    setupNode() {
        super.setupNode();
        
        //setup event handlers
        let self = this;
        
        let checkbox = this.node.querySelector(".roller-shutter-checkbox");
        checkbox.addEventListener(
            "change", function () {
                self.handleRollerShutterCheckboxChange.call(self, checkbox.checked);
            }
        );
    }

    handleRollerShutterCheckboxChange(checked) {
        let productGroupSelector = this.node.querySelector(".prod-group");
        if (checked) {

            //calculate roller shutter price
            
            this.rollerShutterPrice = 70*this.prodData.quantity + 
                (this.prodData.width+20)*(this.prodData.height+250)/(1000*1000)*150*this.prodData.quantity;
            this.setGroup("B");
            if (this.casing != null) {
                this.casing.setGroup("B");
            }

        } else {
            this.setGroup("A");
            if (this.casing != null) {
                this.casing.setGroup("A");
            }
            this.rollerShutterPrice = 0;
        }
    }

    setupCasingSelector(casingList, productsCollection) {
        let self = this;
        let casingSelector = this.node.querySelector(".casing-selector");
        
        this.addOptions(casingSelector, casingList);

        casingSelector.addEventListener(
            "change", function() {
                self.handleCasingSelection.call(self, productsCollection, casingSelector.value);
            }
        );
    }

    addOptions(casingSelector, casingList) {
        for (let casing of casingList) {
            let option = document.createElement("option");
            option.value = casing.uuid;
            option.textContent = `Posizione ${casing.position}`;
            casingSelector.appendChild(option);
        }
    }

    toggleCasingOptionByValue(value) {
        let option = this.node.querySelector(
            `option[value=${CSS.escape(value)}]`
        );
        option.disabled = !option.disabled;
    }
    

    handleCasingSelection(productsCollection, selectedValue) {
        let event = {name: "casingselectionchange"};


        if(selectedValue != "") {

            this.casing = productsCollection[selectedValue];
            event.casing_uuid = this.casing.uuid;
            // send an event to notify all select elements
            // and delete the entry selected
            this.dispatch(event);

            // change the group of the casing if the fixture
            // linked has rollers shutters
            if(this.rollerShutterPrice > 0) {
                this.casing.setGroup("B");
            }
        
        } else {
            if (this.casing != null) {
                event.casing_uuid = this.casing.uuid;
                this.casing.setGroup("A");
                this.casing = null;
                this.dispatch(event);
            }
        }
    }

    toJson() {
        let jsonObj = super.toJson();

        jsonObj["roller_shutter_price"] = this.rollerShutterPrice;
        
        if (this.casing != null)
            jsonObj["casing"] = this.casing.uuid;
        else
            jsonObj["casing"] = "0";

        return jsonObj;
    }
}