const input = document.getElementById("data-file");
const submitBtn = document.getElementById("submit-1");
const quoteDataContainer = document.getElementById("quote-data-container");
const quoteHeader = document.getElementById("quote-header");
const quoteItemsList = document.getElementById("quote-items-list");

let QUOTE_DATA = {};


submitBtn.addEventListener("click", () => {
  let formData = new FormData();
  formData.append("file", input.files[0]);

  let promise = makeProcessQuoteRequest(formData);

  promise.then(
    (data) => { processQuoteData(data) }
  );

});


async function makeProcessQuoteRequest(formData) {
  let response = await fetch(
    "http://localhost:8888/api/quote_data",
    {
      method: "POST",
      body: formData,
    }
  );

  return await response.json();
}


function processQuoteData(data) {
  console.log(data);
  quoteHeader.textContent = data.title;

  //format header
  QUOTE_DATA["header"] = {
    "commission": data.commission,
    "title": data.title,
  };

  //format products
  for (prod of data.products) {
    let uuid = formatProduct(prod);
    let prodData = QUOTE_DATA[uuid];
    createProductItem(prodData, uuid);
  }


}

function formatProduct(product) {
  const uuid = crypto.randomUUID();

  QUOTE_DATA[uuid] = {
    "height": product.Height,
    "width": product.Width,
    "notes": product.Notes,
    "tot_price": product.Price,
    "product_id": product.ProductId,
    "quantity": product.Quantity,
    "reference": product.Reference,
  };

  return uuid;
}

function createProductItem(product, uuid) {
  let prodElement = new ProductItem(uuid, product);
  prodElement.setupNode();
  prodElement.render(quoteItemsList);
}
